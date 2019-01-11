package v1

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Errors  []string `json:"errors,omitempty"` // application-level error message, for debugging
}

func ErrNotFound() render.Renderer {
	return &ErrResponse{
		Err:            nil,
		HTTPStatusCode: http.StatusNotFound,
		Errors:      []string{ "not found" },
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err: err,
		HTTPStatusCode: http.StatusInternalServerError,
		Errors: []string { err.Error() },
	}
}

func ErrPayloadTooLarge() render.Renderer {
	return &ErrResponse{
		Err: nil,
		HTTPStatusCode: http.StatusRequestEntityTooLarge,
		Errors: []string{"the payload of the request exceeds the maximum upload size"},
	}
}