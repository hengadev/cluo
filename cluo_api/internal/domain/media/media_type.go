package domain

// MediaType represents the type of media file
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	// NOTE: is that the right name ?
	MediaTypeAudio MediaType = "audio"
)
