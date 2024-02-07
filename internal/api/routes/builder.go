package routes

import (
	"net/http"

	"github.com/dngrudin/notes-service/internal/api/handlers/create"
	"github.com/dngrudin/notes-service/internal/api/handlers/delete"
	"github.com/dngrudin/notes-service/internal/api/handlers/get"
	"github.com/dngrudin/notes-service/internal/api/handlers/list"
	"github.com/dngrudin/notes-service/internal/api/handlers/update"
	"github.com/dngrudin/notes-service/internal/api/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Builder struct {
	executor Executable
}

type Executable interface {
	create.NoteCreator
	list.NotesGetter
	get.NoteGetter
	delete.NoteDeleter
	update.NoteUpdater
}

func NewBuilder(executor Executable) Builder {
	return Builder{executor: executor}
}

func (rb *Builder) Build() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(logger.New())
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/notes", func(r chi.Router) {
		r.Get("/", list.New(rb.executor))
		r.Post("/", create.New(rb.executor))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", get.New(rb.executor))
			r.Put("/", update.New(rb.executor))
			r.Delete("/", delete.New(rb.executor))
		})
	})

	return r
}
