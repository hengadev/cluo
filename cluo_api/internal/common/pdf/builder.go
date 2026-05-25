package pdf

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

// buildPDF creates a minimal but valid PDF containing the rendered text.
// It parses the supported HTML tags and produces formatted text in the PDF.
func buildPDF(html string) ([]byte, error) {
	lines := extractLines(html)
	return renderPDF(lines)
}

// textLine represents a line of formatted text in the PDF.
type textLine struct {
	text   string
	size   float64
	indent float64
}

// extractLines parses simplified HTML into text lines with formatting.
func extractLines(html string) []textLine {
	var lines []textLine
	s := html

	for len(s) > 0 {
		// Find the next opening tag we understand
		nextTag := findNextTag(s)
		if nextTag == "" {
			// No more tags, add remaining text
			if text := cleanInlineTags(strings.TrimSpace(s)); text != "" {
				lines = append(lines, textLine{text: text, size: 12, indent: 0})
			}
			break
		}

		tagPos := strings.Index(s, "<"+nextTag)
		if tagPos < 0 {
			break
		}

		// Add any text before the tag (as plain paragraph)
		if tagPos > 0 {
			if text := cleanInlineTags(strings.TrimSpace(s[:tagPos])); text != "" {
				lines = append(lines, textLine{text: text, size: 12, indent: 0})
			}
		}

		closeTag := "</" + nextTag
		tagContent, after, found := extractBlock(s[tagPos:], nextTag, closeTag)
		if !found {
			break
		}

		text := cleanInlineTags(strings.TrimSpace(tagContent))
		if text != "" {
			switch nextTag {
			case "h1":
				lines = append(lines, textLine{text: text, size: 22, indent: 0})
			case "h2":
				lines = append(lines, textLine{text: text, size: 18, indent: 0})
			case "h3":
				lines = append(lines, textLine{text: text, size: 14, indent: 0})
			case "blockquote":
				lines = append(lines, textLine{text: text, size: 11, indent: 20})
			case "li":
				lines = append(lines, textLine{text: "• " + text, size: 12, indent: 15})
			default:
				lines = append(lines, textLine{text: text, size: 12, indent: 0})
			}
		}
		s = after
	}

	return lines
}

// findNextTag finds the next recognized block-level tag in the string.
func findNextTag(s string) string {
	tags := []string{"h1", "h2", "h3", "p", "blockquote", "li"}
	earliest := len(s)
	found := ""
	for _, tag := range tags {
		pos := strings.Index(s, "<"+tag)
		// Match exact: <h1> or <h1 but not <h10>
		if pos >= 0 && pos < earliest {
			// Check it's an exact tag match (followed by > or space)
			after := pos + len(tag) + 1
			if after < len(s) {
				ch := s[after]
				if ch == '>' || ch == ' ' {
					earliest = pos
					found = tag
				}
			}
		}
	}
	return found
}

// extractBlock extracts the content between opening <tag...> and closing </tag>,
// returns (content, remaining, found).
func extractBlock(s string, openTag, closeTag string) (string, string, bool) {
	// Find end of opening tag
	gt := strings.Index(s, ">")
	if gt < 0 {
		return "", "", false
	}
	innerStart := gt + 1

	// Find closing tag
	closePos := strings.Index(s[innerStart:], closeTag)
	if closePos < 0 {
		return "", "", false
	}
	content := s[innerStart : innerStart+closePos]

	// Skip past closing tag
	afterClose := innerStart + closePos + len(closeTag)
	endGt := strings.Index(s[afterClose:], ">")
	if endGt < 0 {
		return content, "", true
	}
	remaining := s[afterClose+endGt+1:]

	return content, remaining, true
}

// Strips inline formatting tags and bare <p> wrappers that TipTap emits inside
// block elements (e.g. <li><p>text</p></li>, <blockquote><p>text</p></blockquote>).
var inlineTagRe = regexp.MustCompile(`</?(?:strong|em|u|p)>`)

// cleanInlineTags strips inline formatting tags (strong, em, u) from text.
func cleanInlineTags(s string) string {
	return inlineTagRe.ReplaceAllString(s, "")
}

// ---------------------------------------------------------------------------
// Minimal PDF renderer
// ---------------------------------------------------------------------------

// renderPDF builds a valid PDF 1.4 document from text lines.
func renderPDF(lines []textLine) ([]byte, error) {
	if len(lines) == 0 {
		// Generate a PDF with a blank page
		lines = []textLine{{text: "", size: 12, indent: 0}}
	}

	var buf bytes.Buffer
	objOffsets := make([]int64, 0, 20)

	// Helper to track object offsets
	writeObj := func(objNum int, content string) {
		objOffsets = append(objOffsets, int64(buf.Len()))
		fmt.Fprintf(&buf, "%d 0 obj\n%s\nendobj\n", objNum, content)
	}

	// PDF Header
	fmt.Fprint(&buf, "%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")

	// Object 1: Catalog
	writeObj(1, "<< /Type /Catalog /Pages 2 0 R >>")

	// Object 2: Pages
	writeObj(2, "<< /Type /Pages /Kids [3 0 R] /Count 1 >>")

	// Object 3: Page
	writeObj(3, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>")

	// Build content stream
	var content bytes.Buffer
	// Use Helvetica (WinAnsiEncoding supports French characters well enough for basic use)
	// We'll use Helvetica for normal text and approximate sizes.

	y := 800.0 // Start near top of A4

	for _, line := range lines {
		if y < 50 {
			slog.Warn("pdf: content truncated — document exceeds single page capacity")
			break
		}
		if line.text == "" {
			y -= line.size * 0.5
			continue
		}

		x := 50 + line.indent
		content.WriteString(fmt.Sprintf("BT /F1 %.1f Tf %.1f %.1f Td (%s) Tj ET\n",
			line.size, x, y, pdfEscape(line.text)))
		y -= line.size * 1.5
	}

	contentStr := content.String()

	// Object 4: Content stream
	writeObj(4, fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(contentStr), contentStr))

	// Object 5: Font
	writeObj(5, "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>")

	// Info object (6)
	writeObj(6, fmt.Sprintf("<< /Producer (Cluo PDF Generator) /CreationDate (D:%s) >>",
		time.Now().UTC().Format("20060102150405")))

	// Cross-reference table
	xrefOffset := buf.Len()
	fmt.Fprint(&buf, "xref\n")
	fmt.Fprintf(&buf, "0 %d\n", len(objOffsets)+1)
	fmt.Fprint(&buf, "0000000000 65535 f \n")
	for _, off := range objOffsets {
		fmt.Fprintf(&buf, "%010d 00000 n \n", off)
	}

	// Trailer
	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root 1 0 R /Info 6 0 R >>\n", len(objOffsets)+1)
	fmt.Fprintf(&buf, "startxref\n%d\n%%%%EOF\n", xrefOffset)

	return buf.Bytes(), nil
}

// pdfEscape escapes special characters for PDF string literals.
func pdfEscape(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	for _, r := range s {
		switch r {
		case '\\':
			buf.WriteString("\\\\")
		case '(':
			buf.WriteString("\\(")
		case ')':
			buf.WriteString("\\)")
		case '•':
			buf.WriteString("\\267") // bullet in WinAnsiEncoding
		default:
			if r >= 32 && r <= 255 {
				// WinAnsiEncoding expects single bytes; WriteRune would emit
				// multi-byte UTF-8 for code points > 127, corrupting the output.
				buf.WriteByte(byte(r))
			} else if r > 255 {
				io.WriteString(&buf, fmt.Sprintf("\\%03o", r))
			}
		}
	}
	return buf.String()
}
