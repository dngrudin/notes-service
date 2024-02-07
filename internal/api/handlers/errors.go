package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Error          error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var (
	ErrBadRequest = &ErrResponse{HTTPStatusCode: 400, StatusText: "Bad request"}
	ErrNotFound   = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found"}
	ErrConflict   = &ErrResponse{HTTPStatusCode: 409, StatusText: "Already exists"}

	ErrInternal = &ErrResponse{HTTPStatusCode: 500, StatusText: "Internal server error"}
)
