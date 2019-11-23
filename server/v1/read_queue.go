package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"noonhack/errors"
	"noonhack/respond"

	"github.com/go-chi/chi"
)

func readQueuesHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	queueName := chi.URLParam(r, "queueName")

	queueMu.RLock()
	qinfo, ok := FileQueue[queueName]
	if !ok {
		queueMu.RUnlock()
		return errors.BadRequest("invalid queue name")
	}
	queueMu.RUnlock()
	output, err := readFile(qinfo.filePath)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	respond.OK(w, output, nil)
	return nil
}

func readFile(path string) (map[string]*queueDataInfo, error) {
	queueMu.RLock()
	defer queueMu.RUnlock()

	bytestream, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var output = make(map[string]*queueDataInfo)
	if err := json.Unmarshal(bytestream, &output); err != nil {
		return nil, err
	}

	return output, nil
}
