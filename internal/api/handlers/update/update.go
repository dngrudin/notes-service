package update

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

type NoteUpdater interface {
	UpdateNote(note model.Note) *handlers.ErrResponse
}

type updateNoteRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (cr *updateNoteRequest) Bind(r *http.Request) error {
	return nil
}

func New(noteUpdater NoteUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := slog.Default().With(slog.String("handler", "update"))

		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Error("Bad request", logger.Err(err))
			render.Render(w, r, handlers.ErrBadRequest)
			return
		}

		data := &updateNoteRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Error("Bad request", logger.Err(err))
			render.Render(w, r, handlers.ErrBadRequest)
			return
		}

		errResp := noteUpdater.UpdateNote(model.NewNote(id, data.Title, data.Text))
		if errResp != nil {
			log.Error(errResp.StatusText)
			render.Render(w, r, errResp)
			return
		}

		w.WriteHeader(200)
		w.Write(nil)
	}
}
