package async

import (
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

const (
	TaskNameApiCall   = "apiCall"
	TaskNameApiResult = "apiResult"
)

type apiCallPayload struct {
	Project     *models.Project
	ChainClient *wasm.ChainClient
	Data        []byte
}

func NewApiCallTask(prj *models.Project, chainCli *wasm.ChainClient, data []byte) (*asynq.Task, error) {
	payload, err := json.Marshal(apiCallPayload{
		Project:     prj,
		ChainClient: chainCli,
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
