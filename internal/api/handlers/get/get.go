package get

import (
	"log/slog"
	"net/http"

	"github.com/dngrudin/notes-service/internal/api/handlers"
	"github.com/dngrudin/notes-service/internal/model"
	"github.com/dngrudin/notes-service/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/google/uuid"
)

type NoteGetter interface {
	GetNoteByID(id uuid.UUID) (model.Note, *handlers.ErrResponse)
}

type noteResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Text  string    `json:"text"`
}

func (nr noteResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newNoteResponse(n model.Note) noteResponse {
	return noteResponse{
		ID:    n.ID,
		Title: n.Title,
		Text:  n.Text,
	}
}

func New(noteGetter NoteGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := slog.Default().With(slog.String("handler", "get"))

		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Error("Bad request", logger.Err(err))
			render.Render(w, r, handlers.ErrBadRequest)
			return
		}

		note, errResp := noteGetter.GetNoteByID(id)
		if errResp != nil {
			log.Error(errResp.StatusText)
			render.Render(w, r, errResp)
			return
		}

		render.Render(w, r, newNoteResponse(note))
	}
}
