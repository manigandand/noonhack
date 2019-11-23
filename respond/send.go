package respond

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Response struct contains all the fields needed to respond
// to a particular request
type Response struct {
	StatusCode int
	Data       interface{}
	Headers    map[string]string
}

// SendResponse is a helper function which sends a response with the passed data
func SendResponse(w http.ResponseWriter, statusCode int, data interface{}, headers map[string]string) error {
	return NewResponse(statusCode, data, headers).Send(w)
}

// NewResponse returns a new response object.
func NewResponse(statusCode int, data interface{}, headers map[string]string) *Response {
	return &Response{
		StatusCode: statusCode,
		Data:       data,
		Headers:    headers,
	}
}

// Send sends data encoded to JSON
func (res *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	if res.Headers != nil {
		for key, value := range res.Headers {
			w.Header().Set(key, value)
		}
	}
	w.WriteHeader(res.StatusCode)

	if res.StatusCode != http.StatusNoContent {
		if err := json.NewEncoder(w).Encode(res.Data); err != nil {
			log.Error("respond.send.error: ", err)
			// TODO: handle err, notify developers. Configurable
			// http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// FIXME: remove errors
	return nil
}
