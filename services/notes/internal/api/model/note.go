package model

type CreateNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdateNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
