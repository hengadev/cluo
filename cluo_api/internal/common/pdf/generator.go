package pdf

import (
	"bytes"
	"errors"
	"fmt"
)

var ErrEmptyInput = errors.New("pdf: empty HTML input")

// GenerateFromHTML converts a body HTML string into a PDF byte slice.
// Returns ErrEmptyInput if html is empty.
func GenerateFromHTML(html string) ([]byte, error) {
	if html == "" {
		return nil, ErrEmptyInput
	}

	return generate(html)
}

// IsPDF checks whether data starts with the PDF magic bytes.
func IsPDF(data []byte) bool {
	return bytes.HasPrefix(data, []byte("%PDF"))
}

func generate(html string) ([]byte, error) {
	pdfBytes, err := buildPDF(html)
	if err != nil {
		return nil, fmt.Errorf("pdf: generation failed: %w", err)
	}

	if !IsPDF(pdfBytes) {
		return nil, fmt.Errorf("pdf: generated output is not a valid PDF")
	}

	return pdfBytes, nil
}

// wrapHTMLDocument wraps body HTML into a full HTML document with print CSS.
// Reserved for future use with a full HTML renderer (e.g. chromedp).
func wrapHTMLDocument(bodyHTML string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head>
<meta charset="UTF-8">
<style>
  body {
    font-family: Georgia, 'Times New Roman', serif;
    font-size: 12pt;
    line-height: 1.6;
    color: #1a1a1a;
    max-width: 700px;
    margin: 0 auto;
    padding: 40px 20px;
  }
  h1 { font-size: 22pt; margin-top: 32pt; margin-bottom: 12pt; }
  h2 { font-size: 18pt; margin-top: 24pt; margin-bottom: 10pt; }
  h3 { font-size: 14pt; margin-top: 20pt; margin-bottom: 8pt; }
  p  { margin-bottom: 10pt; }
  blockquote {
    border-left: 3px solid #999;
    padding-left: 16px;
    margin: 16pt 0;
    color: #555;
    font-style: italic;
  }
  ul, ol { margin-left: 20pt; margin-bottom: 10pt; }
  li { margin-bottom: 4pt; }
  strong { font-weight: 700; }
  em { font-style: italic; }
  u { text-decoration: underline; }
</style>
</head>
<body>
%s
</body>
</html>`, bodyHTML)
}
