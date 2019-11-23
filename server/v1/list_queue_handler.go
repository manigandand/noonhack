package v1

import (
	"net/http"
	"noonhack/errors"
	"noonhack/respond"
)

func listQueuesHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	var queues []string
	for name := range FileQueue {
		queues = append(queues, name)
	}

	respond.OK(w, queues, nil)
	return nil
}
