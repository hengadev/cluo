package pdf

import (
	"testing"
)

func TestGenerateFromHTML_ValidInput(t *testing.T) {
	html := `<h1>Title</h1><p>Hello world</p>`
	pdfBytes, err := GenerateFromHTML(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF output")
	}
	if !IsPDF(pdfBytes) {
		t.Fatal("output does not start with %PDF magic bytes")
	}
}

func TestGenerateFromHTML_EmptyInput(t *testing.T) {
	_, err := GenerateFromHTML("")
	if err == nil {
		t.Fatal("expected error for empty HTML input")
	}
	if err != ErrEmptyInput {
		t.Fatalf("expected ErrEmptyInput, got: %v", err)
	}
}

func TestGenerateFromHTML_RapportHTML(t *testing.T) {
	// Simulates the output of TipTap → HTML conversion from issue 002
	html := `<h1>Rapport d'investigation</h1>` +
		`<p>Le <strong>suspect</strong> a été observé.</p>` +
		`<blockquote><p>Note importante</p></blockquote>` +
		`<ul><li><p>Point A</p></li><li><p>Point B</p></li></ul>` +
		`<ol><li><p>Étape 1</p></li><li><p>Étape 2</p></li></ol>` +
		`<h2>Conclusion</h2>` +
		`<p><em>Fin du rapport.</em></p>`

	pdfBytes, err := GenerateFromHTML(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF output")
	}
	if !IsPDF(pdfBytes) {
		t.Fatal("output does not start with %PDF magic bytes")
	}
}

func TestIsPDF_ValidHeader(t *testing.T) {
	data := []byte("%PDF-1.4 some content here")
	if !IsPDF(data) {
		t.Fatal("expected IsPDF to return true for data starting with %PDF")
	}
}

func TestIsPDF_InvalidHeader(t *testing.T) {
	data := []byte("<html>not a pdf</html>")
	if IsPDF(data) {
		t.Fatal("expected IsPDF to return false for non-PDF data")
	}
}

func TestIsPDF_EmptyData(t *testing.T) {
	if IsPDF(nil) {
		t.Fatal("expected IsPDF to return false for nil data")
	}
	if IsPDF([]byte{}) {
		t.Fatal("expected IsPDF to return false for empty data")
	}
}

func TestGenerateFromHTML_OnlyInlineTags(t *testing.T) {
	html := `<p>Just a <strong>bold</strong> and <em>italic</em> paragraph.</p>`
	pdfBytes, err := GenerateFromHTML(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsPDF(pdfBytes) {
		t.Fatal("output is not a valid PDF")
	}
}

func TestGenerateFromHTML_Headings(t *testing.T) {
	html := `<h1>Heading 1</h1><h2>Heading 2</h2><h3>Heading 3</h3>`
	pdfBytes, err := GenerateFromHTML(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsPDF(pdfBytes) {
		t.Fatal("output is not a valid PDF")
	}
}

func TestGenerateFromHTML_Blockquote(t *testing.T) {
	html := `<blockquote><p>Quoted text</p></blockquote>`
	pdfBytes, err := GenerateFromHTML(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsPDF(pdfBytes) {
		t.Fatal("output is not a valid PDF")
	}
}
