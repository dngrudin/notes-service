package list

import (
	"log/slog"
	"net/http"

	"github.com/dngrudin/notes-service/internal/api/handlers"
	"github.com/dngrudin/notes-service/internal/model"

	"github.com/go-chi/render"

	"github.com/google/uuid"
)

type NotesGetter interface {
	GetNotes() ([]model.Note, *handlers.ErrResponse)
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

func newNotesResponse(notes []model.Note) []render.Renderer {
	list := []render.Renderer{}
	for _, note := range notes {
		nr := newNoteResponse(note)
		list = append(list, nr)
	}
	return list
}

func New(notesGetter NotesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := slog.Default().With(slog.String("handler", "create"))

		notes, errResp := notesGetter.GetNotes()
		if errResp != nil {
			log.Error(errResp.StatusText)
			render.Render(w, r, errResp)
			return
		}

		render.RenderList(w, r, newNotesResponse(notes))
	}
}
