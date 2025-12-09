package domain

// MediaType represents the type of media file
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
)

// IsValid checks if the MediaType is valid
func (mt MediaType) IsValid() bool {
	switch mt {
	case MediaTypeImage, MediaTypeVideo, MediaTypeAudio:
		return true
	default:
		return false
	}
}
