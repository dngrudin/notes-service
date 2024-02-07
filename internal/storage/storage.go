package storage

import (
	"github.com/dngrudin/notes-service/internal/model"

	"github.com/google/uuid"
)

type Storage interface {
	GetNotes() ([]model.Note, error)
	GetNoteByID(id uuid.UUID) (model.Note, error)
	CreateNote(note model.Note) error
	UpdateNote(note model.Note) error
	DeleteNoteByID(id uuid.UUID) error
}
