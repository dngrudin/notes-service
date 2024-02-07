package service

import (
	"errors"

	"github.com/dngrudin/notes-service/internal/api/handlers"
	"github.com/dngrudin/notes-service/internal/model"
	"github.com/dngrudin/notes-service/internal/storage"

	"github.com/google/uuid"
)

type Executer struct {
	storage storage.Storage
}

func NewExecuter(storage storage.Storage) *Executer {
	return &Executer{storage: storage}
}

func (e *Executer) CreateNote(note model.Note) *handlers.ErrResponse {
	if err := e.storage.CreateNote(note); err != nil {
		if errors.Is(err, storage.ErrExists) {
			return handlers.ErrConflict
		} else {
			return handlers.ErrInternal
		}
	}
	return nil
}

func (e *Executer) GetNoteByID(id uuid.UUID) (model.Note, *handlers.ErrResponse) {
	note, err := e.storage.GetNoteByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return model.Note{}, handlers.ErrNotFound
		} else {
			return model.Note{}, handlers.ErrInternal
		}
	}

	return note, nil
}

func (e *Executer) GetNotes() ([]model.Note, *handlers.ErrResponse) {
	notes, err := e.storage.GetNotes()
	if err != nil {
		return nil, handlers.ErrInternal
	}

	return notes, nil
}

func (e *Executer) DeleteNoteByID(id uuid.UUID) *handlers.ErrResponse {
	if err := e.storage.DeleteNoteByID(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return handlers.ErrNotFound
		} else {
			return handlers.ErrInternal
		}
	}

	return nil
}

func (e *Executer) UpdateNote(note model.Note) *handlers.ErrResponse {
	if err := e.storage.UpdateNote(note); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return handlers.ErrNotFound
		} else {
			return handlers.ErrInternal
		}
	}

	return nil
}
