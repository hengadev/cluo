// Package tiptap converts TipTap JSON documents to sanitized HTML strings.
//
// Supported node types:
//   - doc (root container)
//   - heading (h1–h3, controlled by attrs.level)
//   - paragraph
//   - text (with optional marks: bold, italic, underline)
//   - orderedList / bulletList
//   - listItem
//   - blockquote
//
// Unknown node types degrade gracefully: their children are rendered inline
// (plain text) without a wrapping element.
package tiptap

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// TipTap JSON model (only the fields we care about)
// ---------------------------------------------------------------------------

type tipTapDoc struct {
	Type    string       `json:"type"`
	Content []tipTapNode `json:"content"`
}

type tipTapNode struct {
	Type    string            `json:"type"`
	Attrs   map[string]any    `json:"attrs"`
	Content []tipTapNode      `json:"content"`
	Text    string            `json:"text"`
	Marks   []tipTapMark      `json:"marks"`
}

type tipTapMark struct {
	Type string         `json:"type"`
	Attrs map[string]any `json:"attrs"`
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// ToHTML converts a raw TipTap JSON byte slice into a sanitized HTML string.
// An empty or nil input returns ("", nil).
func ToHTML(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	var doc tipTapDoc
	if err := json.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("tiptap: invalid JSON: %w", err)
	}

	if doc.Type != "doc" {
		return "", fmt.Errorf("tiptap: expected root node type 'doc', got %q", doc.Type)
	}

	var sb strings.Builder
	for _, child := range doc.Content {
		renderNode(&sb, child)
	}

	return sb.String(), nil
}

// ---------------------------------------------------------------------------
// Internal rendering
// ---------------------------------------------------------------------------

func renderNode(sb *strings.Builder, n tipTapNode) {
	switch n.Type {

	// ---- Block nodes ----
	case "heading":
		level := attrInt(n.Attrs, "level", 1)
		if level < 1 {
			level = 1
		}
		if level > 3 {
			level = 3
		}
		sb.WriteString(fmt.Sprintf("<h%d>", level))
		renderChildren(sb, n.Content)
		sb.WriteString(fmt.Sprintf("</h%d>", level))

	case "paragraph":
		sb.WriteString("<p>")
		renderChildren(sb, n.Content)
		sb.WriteString("</p>")

	case "blockquote":
		sb.WriteString("<blockquote>")
		renderChildren(sb, n.Content)
		sb.WriteString("</blockquote>")

	case "orderedList":
		sb.WriteString("<ol>")
		renderChildren(sb, n.Content)
		sb.WriteString("</ol>")

	case "bulletList":
		sb.WriteString("<ul>")
		renderChildren(sb, n.Content)
		sb.WriteString("</ul>")

	case "listItem":
		sb.WriteString("<li>")
		renderChildren(sb, n.Content)
		sb.WriteString("</li>")

	// ---- Inline nodes ----
	case "text":
		renderText(sb, n)

	// ---- Unknown: degrade gracefully, render children inline ----
	default:
		renderChildren(sb, n.Content)
	}
}

func renderChildren(sb *strings.Builder, children []tipTapNode) {
	for _, child := range children {
		renderNode(sb, child)
	}
}

func renderText(sb *strings.Builder, n tipTapNode) {
	text := escapeHTML(n.Text)
	if text == "" {
		return
	}

	// Apply marks from inside-out: bold → italic → underline (order doesn't
	// matter for HTML correctness, but we wrap consistently).
	for _, mark := range n.Marks {
		switch mark.Type {
		case "bold":
			text = "<strong>" + text + "</strong>"
		case "italic":
			text = "<em>" + text + "</em>"
		case "underline":
			text = "<u>" + text + "</u>"
		}
	}

	sb.WriteString(text)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// attrInt extracts an integer from a node's attrs map, or returns def.
func attrInt(attrs map[string]any, key string, def int) int {
	if attrs == nil {
		return def
	}
	v, ok := attrs[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return def
	}
}

// escapeHTML sanitizes text for safe HTML insertion.
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
