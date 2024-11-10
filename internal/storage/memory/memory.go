package memory

import (
	"sync"

	"github.com/dngrudin/notes-service/internal/model"
	"github.com/dngrudin/notes-service/internal/storage"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	notes map[uuid.UUID]model.Note
	mu    sync.RWMutex
}

func New() (storage.Storage, error) {
	return &MemoryStorage{notes: map[uuid.UUID]model.Note{}}, nil
}

func (s *MemoryStorage) GetNotes() ([]model.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var notes []model.Note
	for _, n := range s.notes {
		notes = append(notes, n)
	}
	return notes, nil
}

func (s *MemoryStorage) GetNoteByID(id uuid.UUID) (model.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n, ok := s.notes[id]
	if !ok {
		return model.Note{}, storage.ErrNotFound
	}

	return n, nil
}

func (s *MemoryStorage) CreateNote(note model.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notes[note.ID]; ok {
		return storage.ErrExists
	}

	s.notes[note.ID] = note
	return nil
}

func (s *MemoryStorage) UpdateNote(note model.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notes[note.ID]; ok {
		s.notes[note.ID] = note
		return nil
	}

	return storage.ErrNotFound
}

func (s *MemoryStorage) DeleteNoteByID(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notes[id]; ok {
		delete(s.notes, id)
		return nil
	}

	return storage.ErrNotFound
}
