package api

import (
	"encoding/json"
	"net/http"

	"example.com/jobqueue/domain"
	"example.com/jobqueue/jobdefs"
	"example.com/jobqueue/processor"
	"example.com/jobqueue/storage"
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
	jobs, err := h.jp.Storage.List()
	if err != nil {
		c.JSON(503, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H {
		"message": "Jobs fetched successfully",
		"jobs": domain.ToJobResponseList(jobs),
	})
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
		"message": "Job created successfully",
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
		"message": "Job fetched successfully",
		"job": domain.ToJobResponse(result)})
}

func (h *JobHandler) CancelJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Job ID is required"})
		return
	}

	err := h.jp.CancelJob(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Job cancellation requested"})
}

func NewJobHandler(jp *processor.JobProcessor) *JobHandler {
	return &JobHandler{jp: jp}
}

type MetricsHandler struct {
	s storage.Storage
}

func NewMetricsHandler(s storage.Storage) *MetricsHandler {
	return &MetricsHandler{s: s}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	// Implement metrics retrieval logic here
	c.JSON(200, gin.H{"message": "Metrics endpoint - to be implemented"})
}
