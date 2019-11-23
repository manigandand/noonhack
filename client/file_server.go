package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"noonhack/errors"
	"noonhack/types"
	"time"
)

type listQueueResponse struct {
	Data []string    `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

type ListQueueDataResponse struct {
	Data []struct {
		Time        time.Time `json:"time"`
		QueueName   string    `json:"queue_name"`
		ServiceName string    `json:"service_name"`
		Data        []byte    `json:"data"`
	} `json:"data"`
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
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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
	if res.StatusCode > 299 && f.client.ToRetry {
		// Retry based on error context
		reqID := time.Now().UnixNano()
		f.client.RetryMU.Lock()
		f.client.RetryCount++
		f.client.RetryStat[reqID] = time.Now()
		f.client.RetryMU.Unlock()

		f.retryPush(reqID, url, body)
	} else {
		whathappened = response["meta"]
	}

	fmt.Println(whathappened)
	return nil
}

func (f *FileServerQueue) retryPush(reqID int64, url string, body []byte) *errors.AppError {
	f.client.RetryMU.RLock()
	retryCount := f.client.RetryCount
	retryStartAT, ok := f.client.RetryStat[reqID]
	retryStartAT = retryStartAT.Add(f.client.RetryInterval)
	if !ok {
		f.client.RetryMU.RUnlock()
		return errors.InternalServer("invalid req id")
	}
	f.client.RetryMU.RUnlock()

	// 1. check if it reach max count
	if retryCount > f.client.MaxRetry {
		return errors.BadRequest("max retry reached")
	}
	// 2. check if it exceeds the interval
	if !retryStartAT.Before(time.Now()) {
		return errors.InternalServer("retry time out exceeded")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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

	if res.StatusCode == 200 {
		f.client.RetryMU.Lock()
		f.client.RetryCount = 0
		f.client.RetryMU.Unlock()
		return nil
	}

	if res.StatusCode > 299 {
		// Retry again based on error context
		f.client.RetryMU.Lock()
		f.client.RetryCount++
		f.client.RetryMU.Unlock()

		f.retryPush(reqID, url, body)
	}

	return errors.InternalServerStd()
}

// Poll implements the QueueClient interface Poll method
func (f *FileServerQueue) Poll(queueName string) (*ListQueueDataResponse, *errors.AppError) {
	url := fmt.Sprintf("%s/v1/queue/%s", f.client.ServerHost, queueName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	defer res.Body.Close()

	var response ListQueueDataResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return &response, nil
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
	// TODO: validate status code

	return response.Data, nil
}
