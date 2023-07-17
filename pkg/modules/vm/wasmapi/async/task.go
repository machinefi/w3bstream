package async

import (
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/machinefi/w3bstream/pkg/models"
)

const (
	TaskNameApiCall   = "apiCall"
	TaskNameApiResult = "apiResult"
)

type apiCallPayload struct {
	Project *models.Project
	Data    []byte
}

func NewApiCallTask(prj *models.Project, data []byte) (*asynq.Task, error) {
	payload, err := json.Marshal(apiCallPayload{
		Project: prj,
		Data:    data,
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
