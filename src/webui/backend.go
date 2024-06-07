package webui

import (
	"github.com/gin-gonic/gin"

	metricService "mc_reverse_proxy/src/metric/service"
	proxyService "mc_reverse_proxy/src/proxy/service"
	common "mc_reverse_proxy/src/webui/backend/common"
	console "mc_reverse_proxy/src/webui/backend/console/controller"
	metric "mc_reverse_proxy/src/webui/backend/metric/controller"
)

type HTTPBackend struct {
	metricController  common.HTTPController
	consoleController common.HTTPController
}

func (b *HTTPBackend) Config(engine *gin.Engine) {
	b.metricController.Config(engine)
	b.consoleController.Config(engine)
}

func NewHTTPBackend(metricCollector *metricService.MetricService, serverRepo proxyService.ServerRepositoryService) common.HTTPController {
	b := &HTTPBackend{
		metricController:  metric.NewMetricController(metricCollector),
		consoleController: console.NewConsoleController(serverRepo),
	}
	return b
}
