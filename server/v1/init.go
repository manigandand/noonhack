package v1

import (
	"net/http"
	"noonhack/errors"
	"noonhack/respond"

	"github.com/go-chi/chi"
	"github.com/gorilla/context"
	"github.com/labstack/gommon/log"
)

// FileQueue a global map that holds all the available queues in the server
// TODO: value might be a struct that holds all the client service information
// as well the file informations for the service
var FileQueue map[string]interface{}

// InitQueue populates static queues for now
func InitQueue() {
	FileQueue = map[string]interface{}{
		"queue_a": "TOBE ADDED",
		"queue_b": "TOBE ADDED",
		"queue_c": "TOBE ADDED",
	}
}

// Init initializes all the v1 routes
func Init(r chi.Router) {
	r.Method(http.MethodGet, "/queue", Handler(listQueuesHandler))
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
