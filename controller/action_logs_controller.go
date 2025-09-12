package controller

import (
	"ecommerce/service"
	_ "ecommerce/entity"
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

// GetLogs godoc
// @Summary      Get all action logs
// @Description  Return list of all action logs (admin only)
// @Tags         Action Logs
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   entity.ActionLog
// @Failure      500  {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch logs, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/logs [get]
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

// GetLogByID godoc
// @Summary      Get log by ID
// @Description  Return a single log by its ID (admin only)
// @Tags         Action Logs
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Log ID (UUID)"
// @Success      200  {object}  entity.ActionLog
// @Failure      400  {object}  map[string]interface{}
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid log ID format",
//   "detail": "some error message"
// }
// @Failure      404  {object}  map[string]interface{}
// @Example 404 {json} Error Example:
// {
//   "message": "No log found",
//   "detail": "some error message"
// }
// @Router       /admin/logs/{id} [get]
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

// ReportSelling godoc
// @Summary      Sales report
// @Description  Get best-selling or least-selling products (admin only)
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        type   query     string  false  "Report type (best/least)"  Enums(best,least)  default(best)
// @Param        limit  query     int     false  "Limit number of results"   default(5)
// @Success      200    {array}   map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch sales report, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/reports/selling [get]
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

// ReportStock godoc
// @Summary      Stock report
// @Description  Get low-stock or high-stock products (admin only)
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        type   query     string  false  "Report type (low/high)"  Enums(low,high)  default(low)
// @Param        limit  query     int     false  "Limit number of results" default(5)
// @Success      200    {array}   entity.Product
// @Failure      500    {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch stock report, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/reports/stock [get]
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
