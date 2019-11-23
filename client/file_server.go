package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"noonhack/errors"
)

type listQueueResponse struct {
	Data []string    `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

// FileServerQueue ...
type FileServerQueue struct {
	client *Client
}

// Push implements the QueueClient interface Push method
func (f *FileServerQueue) Push(queueName string) {

}

// Poll implements the QueueClient interface Poll method
func (f *FileServerQueue) Poll(queueName string) {

}

// ListQueue implements the QueueClient interface ListQueue method
// returns all the available queues from the server
func (f *FileServerQueue) ListQueue() ([]string, *errors.AppError) {
	url := fmt.Sprintf("%s/v1/queue", f.client.ServerHost)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []string{}, errors.InternalServer(err.Error())
	}
	// TODO: add any headers which required to identify the service/consumer
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return []string{}, errors.InternalServer(err.Error())
	}
	defer res.Body.Close()

	var response listQueueResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return []string{}, errors.InternalServer(err.Error())
	}

	return response.Data, nil
}
