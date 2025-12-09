package mediaHandler

const (
	MediaBasePath = "/media"

	// Media CRUD operations
	UploadMediaEndpoint     = MediaBasePath
	GetMediaByIDEndpoint    = MediaBasePath + "/{id}"
	UpdateMediaEndpoint     = MediaBasePath + "/{id}"
	DeleteMediaEndpoint     = MediaBasePath + "/{id}"
	ListMediaByCaseEndpoint = "/case/{caseId}/media"
)
