package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"noonhack/errors"
	"noonhack/types"
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
func (f *FileServerQueue) Push(queueName string, input interface{}) *errors.AppError {
	url := fmt.Sprintf("%s/v1/queue/%s", f.client.ServerHost, queueName)
	data := &types.QueueInput{
		ServiceName: f.client.ServiceName,
		Data:        input,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return errors.BadRequest(err.Error())
	}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(body))
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return errors.InternalServer(err.Error())
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return errors.InternalServer(err.Error())
	}

	whathappened := response["data"]
	if res.StatusCode > 299 {
		whathappened = response["meta"]
	}

	fmt.Println(whathappened)

	return nil
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
