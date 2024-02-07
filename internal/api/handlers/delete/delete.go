package delete

import (
	"log/slog"
	"net/http"

	"github.com/dngrudin/notes-service/internal/api/handlers"
	"github.com/dngrudin/notes-service/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/google/uuid"
)

type NoteDeleter interface {
	DeleteNoteByID(id uuid.UUID) *handlers.ErrResponse
}

func New(noteDeleter NoteDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := slog.Default().With(slog.String("handler", "delete"))

		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Error("Bad request", logger.Err(err))
			render.Render(w, r, handlers.ErrBadRequest)
			return
		}

		errResp := noteDeleter.DeleteNoteByID(id)
		if errResp != nil {
			log.Error(errResp.StatusText)
			render.Render(w, r, errResp)
			return
		}

		w.WriteHeader(200)
		w.Write(nil)
	}
}
