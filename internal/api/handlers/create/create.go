package create

import (
	"log/slog"
	"net/http"

	"github.com/dngrudin/notes-service/internal/api/handlers"
	"github.com/dngrudin/notes-service/internal/model"
	"github.com/dngrudin/notes-service/pkg/logger"

	"github.com/go-chi/render"

	"github.com/google/uuid"
)

type NoteCreator interface {
	CreateNote(note model.Note) *handlers.ErrResponse
}

type createNoteRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (cr *createNoteRequest) Bind(r *http.Request) error {
	return nil
}

func New(noteCreator NoteCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := slog.Default().With(slog.String("handler", "create"))

		data := &createNoteRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Error("Bad request", logger.Err(err))
			render.Render(w, r, handlers.ErrBadRequest)
			return
		}

		id, err := uuid.Parse(data.ID)
		if err != nil {
			log.Error("Bad request", logger.Err(err))
			render.Render(w, r, handlers.ErrBadRequest)
			return
		}

		errResp := noteCreator.CreateNote(model.NewNote(id, data.Title, data.Text))
		if errResp != nil {
			log.Error(errResp.StatusText)
			render.Render(w, r, errResp)
			return
		}

		w.WriteHeader(200)
		w.Write(nil)
	}
}
