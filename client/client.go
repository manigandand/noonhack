package client

import (
	"noonhack/errors"
	"sync"
	"time"
)

const (
	// SERVERHOST ...
	SERVERHOST = "http://localhost:9090"
)

// QueueClient interface
type QueueClient interface {
	Push(queueName string, input interface{}) *errors.AppError
	Poll(queueName string) (*ListQueueDataResponse, *errors.AppError)
	ListQueue() ([]string, *errors.AppError)
}

// Client holds the file queue system config
// TODO: we can have http.Client that can be configurable by the service/consumer
type Client struct {
	sync.RWMutex
	ServerHost string
	// Service Information
	ServiceName string

	// Retry config
	RetryConfig
}

// RetryConfig holds all the config of the service, retry mechanisam etc..
type RetryConfig struct {
	ToRetry       bool
	MaxRetry      int
	RetryInterval time.Duration
	// holds the state of error retry
	RetryMU    sync.Mutex
	RetryCount int
	RetryStat  map[string]time.Time
}

// NewClient constructs the client object and returns QueueClient interface
// to connect
func NewClient(serviceName string) QueueClient {
	client := &Client{
		ServerHost:  SERVERHOST,
		ServiceName: serviceName,
		RetryConfig: RetryConfig{
			ToRetry:   false,
			RetryStat: make(map[string]time.Time, 0),
		},
	}

	return &FileServerQueue{
		client: client,
	}
}

// NewClientWithRetry constructs the client object with retry config
// and returns QueueClient interface to connect
func NewClientWithRetry(serviceName string, maxRetry int, maxInterval time.Duration) QueueClient {
	client := &Client{
		ServerHost:  SERVERHOST,
		ServiceName: serviceName,
		RetryConfig: RetryConfig{
			ToRetry:       true,
			MaxRetry:      maxRetry,
			RetryInterval: maxInterval,
			RetryStat:     make(map[string]time.Time, 0),
		},
	}

	return &FileServerQueue{
		client: client,
	}
}
