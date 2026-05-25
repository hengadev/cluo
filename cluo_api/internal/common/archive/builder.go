// Package archive assembles zip archives for case dossiers and media bundles.
//
// Two variants are provided:
//   - Full archive: documents (PDFs) + rapport (PDF) + published media files.
//   - Media archive: published media files only.
//
// All methods accept an io.Writer so the caller can stream the result directly
// to an HTTP response body without buffering the entire zip in memory.
package archive

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/doctemplate"
	"github.com/hengadev/cluo_api/internal/common/pdf"
	"github.com/hengadev/cluo_api/internal/common/tiptap"
	"github.com/hengadev/cluo_api/internal/domain/document"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
)

// Dependencies abstracts the data sources needed by the archive builder.
type Dependencies interface {
	// ListDocumentsByCase returns all document summaries for a case.
	ListDocumentsByCase(ctx context.Context, caseID uuid.UUID) ([]document.DocumentSummary, error)
	// GetDocument returns the encrypted document for the given ID and type.
	GetDocument(ctx context.Context, id string, docType document.DocumentType) (document.Documentable, error)
	// GetRapportContent returns the TipTap JSON content for the rapport of the given case.
	// If no rapport exists, it returns (nil, nil).
	GetRapportContent(ctx context.Context, caseID uuid.UUID) ([]byte, error)
	// ListPublishedMediaByCase returns all published media for a case (already decrypted).
	ListPublishedMediaByCase(ctx context.Context, caseID uuid.UUID) ([]*domainMedia.MediaFile, error)
	// DownloadMedia fetches the raw bytes for a media file from object storage.
	DownloadMedia(ctx context.Context, url string) (io.ReadCloser, error)
	// DecryptDocument decrypts an encrypted document into a usable domain object.
	DecryptDocument(ctx context.Context, encDoc document.Documentable) (interface{}, error)
}

// BuildFullArchive writes a zip to w containing:
//   - PDFs for each non-draft document,
//   - the rapport as PDF (if present),
//   - all published media files.
//
// If an individual S3 fetch fails for a media file, the file is skipped and
// the archive continues — the caller receives a valid zip with the files that
// could be retrieved. This "best-effort" approach avoids failing the entire
// download when a single media file is unavailable.
func BuildFullArchive(ctx context.Context, deps Dependencies, caseID uuid.UUID, w io.Writer) (err error) {
	zw := zip.NewWriter(w)
	defer func() {
		if cerr := zw.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("archive: close zip: %w", cerr)
		}
	}()

	if e := addDocumentPDFs(ctx, deps, caseID, zw); e != nil {
		return fmt.Errorf("archive: add documents: %w", e)
	}

	if e := addRapportPDF(ctx, deps, caseID, zw); e != nil {
		return fmt.Errorf("archive: add rapport: %w", e)
	}

	if e := addMediaFiles(ctx, deps, caseID, zw); e != nil {
		return fmt.Errorf("archive: add media: %w", e)
	}

	return nil
}

// BuildMediaArchive writes a zip to w containing only the published media files.
// Uses the same best-effort approach for S3 fetch failures.
func BuildMediaArchive(ctx context.Context, deps Dependencies, caseID uuid.UUID, w io.Writer) (err error) {
	zw := zip.NewWriter(w)
	defer func() {
		if cerr := zw.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("archive: close zip: %w", cerr)
		}
	}()

	if e := addMediaFiles(ctx, deps, caseID, zw); e != nil {
		return fmt.Errorf("archive: add media: %w", e)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func addDocumentPDFs(ctx context.Context, deps Dependencies, caseID uuid.UUID, zw *zip.Writer) error {
	summaries, err := deps.ListDocumentsByCase(ctx, caseID)
	if err != nil {
		return fmt.Errorf("list documents: %w", err)
	}

	for _, summary := range summaries {
		if err := ctx.Err(); err != nil {
			return err
		}

		// Skip draft documents
		if summary.Status == document.DocumentStatusDraft {
			continue
		}

		encDoc, err := deps.GetDocument(ctx, summary.ID.String(), summary.Type)
		if err != nil {
			continue // skip documents that can't be fetched
		}

		decrypted, err := deps.DecryptDocument(ctx, encDoc)
		if err != nil {
			continue // skip documents that can't be decrypted
		}

		html, err := doctemplate.RenderDocument(decrypted)
		if err != nil {
			continue // skip documents that can't be rendered
		}

		pdfBytes, err := pdf.GenerateFromHTML(html)
		if err != nil {
			continue // skip documents whose PDF generation fails
		}

		entryName := fmt.Sprintf("documents/%s.pdf", summary.Type)
		// Include the document ref for a friendlier name when available.
		if summary.DocumentRef != "" {
			entryName = fmt.Sprintf("documents/%s-%s.pdf", summary.DocumentRef, summary.Type)
		}

		f, err := zw.Create(entryName)
		if err != nil {
			return fmt.Errorf("create zip entry %q: %w", entryName, err)
		}

		if _, err := f.Write(pdfBytes); err != nil {
			return fmt.Errorf("write zip entry %q: %w", entryName, err)
		}
	}

	return nil
}

func addRapportPDF(ctx context.Context, deps Dependencies, caseID uuid.UUID, zw *zip.Writer) error {
	content, err := deps.GetRapportContent(ctx, caseID)
	if err != nil {
		return nil // no rapport is not an error
	}
	if len(content) == 0 {
		return nil
	}

	html, err := tiptap.ToHTML(content)
	if err != nil {
		return fmt.Errorf("convert rapport to HTML: %w", err)
	}

	if html == "" {
		return nil
	}

	pdfBytes, err := pdf.GenerateFromHTML(html)
	if err != nil {
		return fmt.Errorf("generate rapport PDF: %w", err)
	}

	f, err := zw.Create("rapport.pdf")
	if err != nil {
		return fmt.Errorf("create zip entry %q: %w", "rapport.pdf", err)
	}

	if _, err := f.Write(pdfBytes); err != nil {
		return fmt.Errorf("write zip entry %q: %w", "rapport.pdf", err)
	}

	return nil
}

func addMediaFiles(ctx context.Context, deps Dependencies, caseID uuid.UUID, zw *zip.Writer) error {
	mediaList, err := deps.ListPublishedMediaByCase(ctx, caseID)
	if err != nil {
		return fmt.Errorf("list media: %w", err)
	}

	seen := make(map[string]int)
	for _, media := range mediaList {
		if err := ctx.Err(); err != nil {
			return err
		}

		rc, err := deps.DownloadMedia(ctx, media.URL)
		if err != nil {
			// Best-effort: skip files that can't be downloaded.
			continue
		}

		entryName := uniqueMediaPath(seen, media.FileName)
		f, err := zw.Create(entryName)
		if err != nil {
			rc.Close()
			return fmt.Errorf("create zip entry %q: %w", entryName, err)
		}

		_, copyErr := io.Copy(f, rc)
		rc.Close()
		if copyErr != nil {
			// Best-effort: skip files that failed mid-stream.
			continue
		}
	}

	return nil
}

// uniqueMediaPath returns a deduplicated, safe path inside the media/ directory.
// When two files share the same base name the second gets a numeric suffix (_2, _3, …).
func uniqueMediaPath(seen map[string]int, fileName string) string {
	base := filepath.Base(fileName)
	base = strings.ReplaceAll(base, " ", "_")
	if base == "" || base == "." || base == ".." {
		base = "unnamed"
	}
	canonical := "media/" + base
	seen[canonical]++
	if seen[canonical] == 1 {
		return canonical
	}
	ext := filepath.Ext(canonical)
	stem := strings.TrimSuffix(canonical, ext)
	return fmt.Sprintf("%s_%d%s", stem, seen[canonical], ext)
}
