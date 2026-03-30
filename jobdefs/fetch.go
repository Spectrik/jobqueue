package jobdefs

import (
	"context"
	"fmt"
	"time"
)
type FetchJobData struct {
    URL string `json:"url" binding:"required,url"`
}

type FetchJob struct {
	data FetchJobData
}

func (j FetchJob) Execute(ctx context.Context) (string, error) {
	select {
		case <-time.After(60 * time.Second):
			return "Fetched.", nil
		case <-ctx.Done():
			fmt.Print(ctx.Deadline())
			fmt.Println("TIMEOUT!")
			return "nil", ctx.Err()
	}
}

func (j FetchJob) Type() string {
	return "fetch"
}
