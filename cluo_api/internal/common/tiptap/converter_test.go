package tiptap

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// mustJSON is a tiny helper to write inline JSON test fixtures.
func mustJSON(t *testing.T, jsonStr string) []byte {
	t.Helper()
	return []byte(jsonStr)
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestToHTML_EmptyInput(t *testing.T) {
	got, err := ToHTML(nil)
	if err != nil {
		t.Fatalf("nil input: unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("nil input: expected empty string, got %q", got)
	}

	got, err = ToHTML([]byte{})
	if err != nil {
		t.Fatalf("empty slice: unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("empty slice: expected empty string, got %q", got)
	}
}

func TestToHTML_EmptyDocument(t *testing.T) {
	input := `{"type":"doc","content":[]}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestToHTML_InvalidJSON(t *testing.T) {
	_, err := ToHTML([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestToHTML_WrongRootType(t *testing.T) {
	input := `{"type":"paragraph","content":[]}`
	_, err := ToHTML(mustJSON(t, input))
	if err == nil {
		t.Fatal("expected error for wrong root type")
	}
}

func TestToHTML_Paragraph(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "paragraph", "content": [
				{"type": "text", "text": "Hello world"}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<p>Hello world</p>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_Headings(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "heading", "attrs": {"level": 1}, "content": [{"type": "text", "text": "H1"}]},
			{"type": "heading", "attrs": {"level": 2}, "content": [{"type": "text", "text": "H2"}]},
			{"type": "heading", "attrs": {"level": 3}, "content": [{"type": "text", "text": "H3"}]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<h1>H1</h1><h2>H2</h2><h3>H3</h3>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_HeadingLevelClamped(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "heading", "attrs": {"level": 0}, "content": [{"type": "text", "text": "Low"}]},
			{"type": "heading", "attrs": {"level": 5}, "content": [{"type": "text", "text": "High"}]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<h1>Low</h1><h3>High</h3>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_BoldItalicUnderline(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "paragraph", "content": [
				{"type": "text", "text": "normal "},
				{"type": "text", "marks": [{"type": "bold"}], "text": "bold"},
				{"type": "text", "text": " "},
				{"type": "text", "marks": [{"type": "italic"}], "text": "italic"},
				{"type": "text", "text": " "},
				{"type": "text", "marks": [{"type": "underline"}], "text": "underlined"}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<p>normal <strong>bold</strong> <em>italic</em> <u>underlined</u></p>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_OrderedList(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "orderedList", "content": [
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "First"}]}]},
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Second"}]}]}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<ol><li><p>First</p></li><li><p>Second</p></li></ol>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_BulletList(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "bulletList", "content": [
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "A"}]}]},
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "B"}]}]}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<ul><li><p>A</p></li><li><p>B</p></li></ul>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_Blockquote(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "blockquote", "content": [
				{"type": "paragraph", "content": [{"type": "text", "text": "Quoted text"}]}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<blockquote><p>Quoted text</p></blockquote>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_UnknownNodeDegradesGracefully(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "customWidget", "content": [
				{"type": "text", "text": "fallback text"}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Unknown node renders its children inline (plain text, no wrapper)
	want := "fallback text"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_DeepNesting(t *testing.T) {
	// blockquote > bulletList > listItem > paragraph > bold text
	input := `{
		"type": "doc",
		"content": [
			{"type": "blockquote", "content": [
				{"type": "bulletList", "content": [
					{"type": "listItem", "content": [
						{"type": "paragraph", "content": [
							{"type": "text", "marks": [{"type": "bold"}], "text": "Deep"}
						]}
					]}
				]}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<blockquote><ul><li><p><strong>Deep</strong></p></li></ul></blockquote>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_HTMLInTextIsEscaped(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "paragraph", "content": [
				{"type": "text", "text": "<script>alert('xss')</script>"}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_MultipleMarksOnSameText(t *testing.T) {
	input := `{
		"type": "doc",
		"content": [
			{"type": "paragraph", "content": [
				{"type": "text", "marks": [{"type": "bold"}, {"type": "italic"}, {"type": "underline"}], "text": "All"}
			]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<p><u><em><strong>All</strong></em></u></p>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestToHTML_ComplexDocument(t *testing.T) {
	// A realistic document with all supported node types
	input := `{
		"type": "doc",
		"content": [
			{"type": "heading", "attrs": {"level": 1}, "content": [{"type": "text", "text": "Rapport d'investigation"}]},
			{"type": "paragraph", "content": [{"type": "text", "text": "Le "} , {"type": "text", "marks": [{"type": "bold"}], "text": "suspect"}, {"type": "text", "text": " a été observé."}]},
			{"type": "blockquote", "content": [
				{"type": "paragraph", "content": [{"type": "text", "text": "Note importante"}]}
			]},
			{"type": "bulletList", "content": [
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Point A"}]}]},
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Point B"}]}]}
			]},
			{"type": "orderedList", "content": [
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Étape 1"}]}]},
				{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Étape 2"}]}]}
			]},
			{"type": "heading", "attrs": {"level": 2}, "content": [{"type": "text", "text": "Conclusion"}]},
			{"type": "paragraph", "content": [{"type": "text", "marks": [{"type": "italic"}], "text": "Fin du rapport."}]}
		]
	}`
	got, err := ToHTML(mustJSON(t, input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `<h1>Rapport d&#39;investigation</h1>` +
		`<p>Le <strong>suspect</strong> a été observé.</p>` +
		`<blockquote><p>Note importante</p></blockquote>` +
		`<ul><li><p>Point A</p></li><li><p>Point B</p></li></ul>` +
		`<ol><li><p>Étape 1</p></li><li><p>Étape 2</p></li></ol>` +
		`<h2>Conclusion</h2>` +
		`<p><em>Fin du rapport.</em></p>`
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
