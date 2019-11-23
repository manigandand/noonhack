package client

import "noonhack/errors"

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
	return []string{}, nil
}
