package ai

// TextTransformationType represents the type of text transformation to apply.
type TextTransformationType string

const (
	// TransformationTypeReword rewrites the text with different wording while preserving meaning.
	TransformationTypeReword TextTransformationType = "reword"
	// TransformationTypeSummarize creates a concise summary of the text.
	TransformationTypeSummarize TextTransformationType = "summarize"
	// TransformationTypeFormalize converts informal text to formal language.
	TransformationTypeFormalize TextTransformationType = "formalize"
	// TransformationTypeClarify simplifies complex text for better understanding.
	TransformationTypeClarify TextTransformationType = "clarify"
)

// Prompt returns the system prompt template for the transformation type.
func (t TextTransformationType) Prompt() string {
	switch t {
	case TransformationTypeReword:
		return `You are a text transformation assistant. Rewrite the provided text using different wording while preserving the original meaning, tone, and key details. Output only the transformed text without explanations or formatting.`
	case TransformationTypeSummarize:
		return `You are a text transformation assistant. Create a concise summary of the provided text, capturing the main points and key details. Output only the summary without explanations or formatting.`
	case TransformationTypeFormalize:
		return `You are a text transformation assistant. Convert the provided informal text to formal, professional language while preserving the original meaning. Output only the transformed text without explanations or formatting.`
	case TransformationTypeClarify:
		return `You are a text transformation assistant. Simplify the provided complex text to make it clearer and easier to understand, while preserving the original meaning. Output only the transformed text without explanations or formatting.`
	default:
		return `You are a text transformation assistant. Transform the provided text as requested. Output only the transformed text without explanations or formatting.`
	}
}

// Valid returns true if the transformation type is valid.
func (t TextTransformationType) Valid() bool {
	switch t {
	case TransformationTypeReword,
		TransformationTypeSummarize,
		TransformationTypeFormalize,
		TransformationTypeClarify:
		return true
	default:
		return false
	}
}
