package api

import (
	"encoding/json"
	"net/http"

	"example.com/jobqueue/jobdefs"
	"example.com/jobqueue/processor"
	"github.com/gin-gonic/gin"
)

type JobRequestSchema struct {
    Kind  string `json:"type" binding:"required"`
    Data  json.RawMessage `json:"data" binding:"required"`
}

type JobHandler struct {
	jp *processor.JobProcessor
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	// Implementation for listing jobs

}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var req JobRequestSchema

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	factory, exists := jobdefs.JobBuilders[req.Kind]
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid job type"})
		return
	}

	job, err := factory(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid, err := h.jp.AddJob(job)
	if err != nil {
		c.JSON(503, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
		"uuid": uuid,
	})
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	result, err := h.jp.Storage.Get(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
		"status": result.Status,
	})
}

func NewJobHandler(jp *processor.JobProcessor) *JobHandler {
	return &JobHandler{jp: jp}
}
