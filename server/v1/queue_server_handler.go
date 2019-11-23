package v1

import (
	"net/http"
	"noonhack/errors"
	"noonhack/respond"
)

func queueServerHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	respond.OK(w, "file saved succesfully", nil)
	return nil
}
