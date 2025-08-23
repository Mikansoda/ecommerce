package controller

import (
	"ecommerce/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActionLogController struct {
	service service.ActionLogService
}

func NewActionLogController(s service.ActionLogService) *ActionLogController {
	return &ActionLogController{service: s}
}

func (ctl *ActionLogController) GetLogs(c *gin.Context) {
	logs, err := ctl.service.GetLogs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch logs, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (ctl *ActionLogController) GetLogByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid log ID format",
			"detail":  err.Error(),
		})
		return
	}
	log, err := ctl.service.GetLogByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No log found",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, log)
}

// combination of best-selling & least-selling report
func (ctl *ActionLogController) ReportSelling(c *gin.Context) {
	// default query best
	reportType := c.DefaultQuery("type", "best") // best/least
	limit := 5
	if val := c.Query("limit"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	result, err := ctl.service.ReportSelling(c, reportType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch sales report, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// combination of low-stock & high-stock report
func (ctl *ActionLogController) ReportStock(c *gin.Context) {
	// default low
	reportType := c.DefaultQuery("type", "low") // low/high
	limit := 5
	if val := c.Query("limit"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	result, err := ctl.service.ReportStock(c, reportType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch stock report, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}
