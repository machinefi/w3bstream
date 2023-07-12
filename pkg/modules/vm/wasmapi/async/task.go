package async

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskNameApiCall   = "apiCall"
	TaskNameApiResult = "apiResult"
)

type apiCallPayload struct {
	ProjectName string
	Data        []byte
}

func NewApiCallTask(projectName string, data []byte) (*asynq.Task, error) {
	payload, err := json.Marshal(apiCallPayload{
		ProjectName: projectName,
		Data:        data,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskNameApiCall, payload), nil
}

type apiResultPayload struct {
	ProjectName string
	EventType   string
	Data        []byte
}

func newApiResultTask(projectName, eventType string, data []byte) (*asynq.Task, error) {
	payload, err := json.Marshal(apiResultPayload{
		ProjectName: projectName,
		EventType:   eventType,
		Data:        data,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskNameApiResult, payload), nil
}
