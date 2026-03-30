package jobdefs

import (
	"encoding/json"

	"example.com/jobqueue/domain"
	"github.com/gin-gonic/gin/binding"
)

func bindJobJSON[T any](raw json.RawMessage) (*T, error) {
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}

	if err := binding.Validator.ValidateStruct(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

var JobBuilders = map[string]domain.JobBuilder{
	"fetch": buildFetchJob,
}

func buildFetchJob(payload json.RawMessage) (domain.Job, error) {
	req, err := bindJobJSON[FetchJobData](payload)

	if err != nil {
		return nil, err
	}

	return FetchJob{
		data: *req,
	}, nil
}
