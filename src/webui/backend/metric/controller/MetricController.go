package controller

import (
	metricService "mc_reverse_proxy/src/metric/service"
	common "mc_reverse_proxy/src/webui/backend/common"

	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricController struct {
	MetricCollector *metricService.MetricService
}

func (c *MetricController) Config(router *gin.Engine) {
	r := router.Group("/metric")
	{
		r.GET("/", func(ctx *gin.Context) {
			metric, err := c.MetricCollector.Collect()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
				return
			} else {
				ctx.JSON(http.StatusOK, metric)
			}
		})
	}
}
func NewMetricController(metricCollector *metricService.MetricService) common.HTTPController {
	return &MetricController{MetricCollector: metricCollector}
}
