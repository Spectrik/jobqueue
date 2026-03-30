package main

import (
	"fmt"

	"example.com/jobqueue/api"
	"example.com/jobqueue/processor"
	"example.com/jobqueue/storage"
)

func main() {
	fmt.Println("Hello, World!")
	r := api.NewRouter()
	s := storage.NewInMemoryStorage()
	p := processor.NewJobProcessor(100, 5, s)
	h := api.NewJobHandler(p)

	r.POST("/job", h.CreateJob)
	r.DELETE("/job/:id", h.CancelJob)
	r.GET("/job/", h.ListJobs)
	r.GET("/job/:id", h.GetJob)
	r.GET("/metrics", api.NewMetricsHandler(s).GetMetrics)
	r.Run()
}
