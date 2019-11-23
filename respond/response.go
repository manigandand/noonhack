package respond

import (
	"net/http"
	"noonhack/errors"
	"noonhack/types"

	log "github.com/sirupsen/logrus"
)

// 2XX -------------------------------------------------------------------------

// OK is a helper function used to send response data
// with StatusOK status code (200)
func OK(w http.ResponseWriter, data, meta interface{}) error {
	return SendResponse(w, http.StatusOK, WrapPayload(data, meta), nil)
}

// Created is a helper function used to send response data
// with StatusCreated status code (201)
func Created(w http.ResponseWriter, data, meta interface{}) error {
	return SendResponse(w, http.StatusCreated, WrapPayload(data, meta), nil)
}

// 4xx & 5XX -------------------------------------------------------------------

// Fail write the error response
// Common func to send all the error response
func Fail(w http.ResponseWriter, e *errors.AppError) {
	log.Errorf("StatusCode: %d, Error: %s\n DEBUG: %s\n",
		e.Status, e.Error(), e.Debug)
	SendResponse(w, e.Status, WrapPayload(nil, e), nil)
}

// WrapPayload is used to create a generic payload for the data
// and the metadata passed
func WrapPayload(data, meta interface{}) types.JSON {
	x := make(types.JSON)
	if data != nil {
		x["data"] = data
	}

	if meta != nil {
		x["meta"] = meta
	}

	return x
}
