package broker

const TopicNoteCreated = "note.created"

type NoteCreatedEvent struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	NoteID   string `json:"note_id"`
	Title    string `json:"title"`
}
