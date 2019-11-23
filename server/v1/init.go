package v1

import (
	"fmt"
	"log"
	"net/http"
	"noonhack/config"
	"noonhack/errors"
	"noonhack/respond"
	"os"
	"sync"

	"github.com/go-chi/chi"
	"github.com/gorilla/context"
)

type queueInfo struct {
	filePath string
}

var queueMu sync.RWMutex

// FileQueue a global map that holds all the available queues in the server
// TODO: value might be a struct that holds all the client service information
// as well the file informations for the service
var FileQueue = make(map[string]queueInfo)

// InitQueue populates static queues for now
func InitQueue() {
	queues := []string{"queue_a", "queue_b", "queue_c"}

	for _, q := range queues {
		// create a queue file if not
		if path, err := createQueueFile(q); err == nil {
			FileQueue[q] = queueInfo{
				filePath: path,
			}
		}
	}

	fmt.Printf("QUEUE |> \n%+v\n", FileQueue)
}

// FileExists checks the given filepath is valid or not
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createQueueFile(name string) (string, error) {
	filePath := fmt.Sprintf("%s/%s.json", config.QueueFileDir, name)
	if FileExists(filePath) {
		return filePath, nil
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Can't clone default config", err)
		return "", err
	}

	if err = f.Close(); err != nil {
		log.Println(err)
		return "", err
	}

	return filePath, nil
}

// Init initializes all the v1 routes
func Init(r chi.Router) {
	r.Method(http.MethodGet, "/queue", Handler(listQueuesHandler))
	r.Method(http.MethodPost, "/queue/{queueName}", Handler(queueServerHandler))
	r.Method(http.MethodGet, "/queue/{queueName}", Handler(readQueuesHandler))
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
		fmt.Printf("StatusCode: %d, Error: %s\n DEBUG: %s\n",
			err.Status, err.Error(), err.Debug)
		respond.Fail(w, err)
	}
}
