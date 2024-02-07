package model

import "github.com/google/uuid"

type Note struct {
	ID    uuid.UUID
	Title string
	Text  string
}

func NewNote(id uuid.UUID, title string, text string) Note {
	return Note{ID: id, Title: title, Text: text}
}
