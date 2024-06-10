package webui

import (
	"github.com/gin-gonic/gin"

	metricService "mc_reverse_proxy/src/metric/service"
	proxyService "mc_reverse_proxy/src/proxy/service"
	common "mc_reverse_proxy/src/webui/backend/common"
	console "mc_reverse_proxy/src/webui/backend/console/controller"
	metric "mc_reverse_proxy/src/webui/backend/metric/controller"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type HTTPBackend struct {
	metricController  common.HTTPController
	consoleController common.HTTPController

	address string
}

func (b *HTTPBackend) Config(engine *gin.Engine) {
	engine.Use(CORSMiddleware())
	b.metricController.Config(engine)
	b.consoleController.Config(engine)
}

func (b *HTTPBackend) Serve() error {
	engine := gin.Default()
	b.Config(engine)
	return engine.Run(b.address)
}

func NewHTTPBackend(address string, metricCollector *metricService.MetricService, serverRepo proxyService.ServerRepositoryService) *HTTPBackend {
	b := &HTTPBackend{
		metricController:  metric.NewMetricController(metricCollector),
		consoleController: console.NewConsoleController(serverRepo),
		address:           address,
	}
	return b
}
