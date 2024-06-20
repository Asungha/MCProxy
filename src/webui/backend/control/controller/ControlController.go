package controller

import (
	"mc_reverse_proxy/src/control/controlProto"
	service "mc_reverse_proxy/src/control/service"
	"mc_reverse_proxy/src/webui/backend/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControlController struct {
	eventService *service.EventService
}

func (c *ControlController) Config(router *gin.Engine) {
	r := router.Group("/control")
	{
		r.POST("/", func(ctx *gin.Context) {
			c.eventService.Publish("command", service.EventData{CommandData: &controlProto.CommandData{Command: controlProto.CommandEnum_TIMESET, TimesetData: &controlProto.TimesetData{Time: 0}}})
			ctx.Status(http.StatusOK)
		})
	}
}
func NewControlController(eventService *service.EventService) common.HTTPController {
	return &ControlController{eventService: eventService}
}
