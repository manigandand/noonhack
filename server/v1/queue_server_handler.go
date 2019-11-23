package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"noonhack/errors"
	"noonhack/respond"
	"noonhack/types"
	"time"

	"github.com/go-chi/chi"
)

type queueDataInfo struct {
	Time        time.Time   `json:"time"`
	QueueName   string      `json:"queue_name"`
	ServiceName string      `json:"service_name"`
	Data        interface{} `json:"data"`
	DataBytes   []byte      `json:"data_bytes"`
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
		Data:        input.Data,
		DataBytes:   dataBytes,
	}
	if err := writeDataToFile(qinfo.filePath, data); err != nil {
		return errors.InternalServer(err.Error())
	}

	respond.OK(w, "file saved succesfully", nil)
	return nil
}

func writeDataToFile(path string, content *queueDataInfo) error {
	output, err := readFile(path)
	if err != nil {
		return errors.InternalServer(err.Error())
	}
	output[fmt.Sprintf("%+v", time.Now().UnixNano())] = content

	jsonData, err := json.Marshal(output)
	if err != nil {
		return err
	}

	queueMu.Lock()
	defer queueMu.Unlock()

	err = ioutil.WriteFile(path, jsonData, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()

	// if _, err = f.Write(jsonData); err != nil {
	// 	return err
	// }

	return nil
}
