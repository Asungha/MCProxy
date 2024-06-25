package webui

import (
	"github.com/gin-gonic/gin"

	service "mc_reverse_proxy/src/control/service"
	metricService "mc_reverse_proxy/src/metric/service"
	proxyService "mc_reverse_proxy/src/proxy/service"
	"mc_reverse_proxy/src/utils"
	common "mc_reverse_proxy/src/webui/backend/common"
	control "mc_reverse_proxy/src/webui/backend/control/controller"
	metric "mc_reverse_proxy/src/webui/backend/metric/controller"
	serverlist "mc_reverse_proxy/src/webui/backend/serverlist/controller"
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
	Controller []common.HTTPController
	address    string
}

func (b *HTTPBackend) Config(engine *gin.Engine) {
	for _, c := range b.Controller {
		c.Config(engine)
	}
}

func (b *HTTPBackend) Serve() error {
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	b.Config(engine)
	// log.Printf("[HTTP Backend] Start server on %s", b.address)
	utils.FLog.HTTPBackend("Start server on %s", b.address)
	return engine.Run(b.address)
}

func NewHTTPBackend(address string, metricCollector *metricService.MetricService, serverRepo proxyService.ServerRepositoryService, eventService *service.EventService) *HTTPBackend {
	gin.SetMode(gin.ReleaseMode)
	b := &HTTPBackend{
		Controller: []common.HTTPController{
			metric.NewMetricController(metricCollector),
			serverlist.NewServerlistController(serverRepo),
			control.NewControlController(eventService),
		},
		address: address,
	}
	return b
}
