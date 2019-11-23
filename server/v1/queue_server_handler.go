package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"noonhack/errors"
	"noonhack/respond"
	"noonhack/types"
	"os"
	"time"

	"github.com/go-chi/chi"
)

type queueDataInfo struct {
	Time        time.Time `json:"time"`
	QueueName   string    `json:"queue_name"`
	ServiceName string    `json:"service_name"`
	Data        []byte    `json:"data"`
}

func queueServerHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	var input types.QueueInput
	queueName := chi.URLParam(r, "queueName")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return errors.BadRequest(err.Error())
	}

	queueMu.RLock()
	qinfo, ok := FileQueue[queueName]
	if !ok {
		queueMu.RUnlock()
		return errors.BadRequest("invalid queue name")
	}
	queueMu.RUnlock()
	dataBytes, err := json.Marshal(input.Data)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	data := &queueDataInfo{
		Time:        time.Now(),
		QueueName:   queueName,
		ServiceName: input.ServiceName,
		Data:        dataBytes,
	}
	writeDataToFile(qinfo.filePath, data)

	respond.OK(w, "file saved succesfully", nil)
	return nil
}

func writeDataToFile(path string, content *queueDataInfo) error {
	queueMu.Lock()
	defer queueMu.Unlock()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	data := map[int64]*queueDataInfo{
		time.Now().UnixNano(): content,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	s := string(jsonData) + "\n"
	if _, err = f.WriteString(s); err != nil {
		return err
	}

	daa, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Print(string(daa))
	return nil
}
