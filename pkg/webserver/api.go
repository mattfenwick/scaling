package webserver

type UploadDocumentRequest struct {
	Document string
}

type UploadDocumentResponse struct {
	DocumentId string
}

type GetDocumentRequest struct {
	DocumentId string
}

type GetDocumentResponse struct {
	Document *Document
}

type UnsafeGetDocumentsResponse struct {
	Documents map[string]*Document
}
