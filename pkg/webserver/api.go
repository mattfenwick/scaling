package webserver

// create document

type UploadDocumentRequest struct {
	Document string
}

type UploadDocumentResponse struct {
	DocumentId string
}

// fetch single document

type GetDocumentRequest struct {
	DocumentId string
}

type GetDocumentResponse struct {
	Document *Document
}

// fetch all documents

type GetAllDocumentsResponse struct {
	Documents map[string]*Document
}

// find documents

type FindDocumentsRequest struct {
	Key string
}

type FindDocumentsResponseItem struct {
	DocumentId string
	Paths      [][]any
}

type FindDocumentsResponse struct {
	Matches []*FindDocumentsResponseItem
}
