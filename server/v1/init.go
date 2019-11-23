package v1

import (
	"net/http"
	"noonhack/errors"
	"noonhack/respond"

	"github.com/go-chi/chi"
	"github.com/gorilla/context"
	"github.com/labstack/gommon/log"
)

// Init initializes all the v1 routes
func Init(r chi.Router) {
	r.Method(http.MethodPost, "/queue", Handler(queueServerHandler))
}

// API Handler's ---------------------------------------------------------------

// Handler custom api handler help us to handle all the errors in one place
type Handler func(w http.ResponseWriter, r *http.Request) *errors.AppError

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	// clear gorilla context
	defer context.Clear(r)
	if err != nil {
		// APP Level Error
		// TODO: handle 5XX, notify developers. Configurable
		log.Errorf("StatusCode: %d, Error: %s\n DEBUG: %s\n",
			err.Status, err.Error(), err.Debug)
		respond.Fail(w, err)
	}
}
