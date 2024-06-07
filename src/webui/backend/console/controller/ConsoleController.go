package controller

import (
	"net/http"

	proxyService "mc_reverse_proxy/src/proxy/service"

	"log"

	"github.com/gin-gonic/gin"
)

type UpdateServerRecord struct {
	ID       int    `json:id`
	Hostname string `json:hostname`
	Address  string `json:address`
}

type ConsoleController struct {
	serverRepo proxyService.ServerRepositoryService
}

func (c *ConsoleController) Config(router *gin.Engine) {
	r := router.Group("/console")
	{
		r.GET("/serverlist", func(ctx *gin.Context) {
			res := []gin.H{}
			list, err := c.serverRepo.(proxyService.ListableRepositoryService).List()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			for _, record := range list {
				res = append(res, gin.H{"id": record.ID, "hostname": record.Hostname, "address": record.Address})
			}
			ctx.JSON(http.StatusOK, res)
		})
		r.PATCH("/serverlist", func(ctx *gin.Context) {
			var req UpdateServerRecord
			if err := ctx.BindJSON(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Printf("%v", req)
			err := c.serverRepo.(proxyService.UpdatableRepositoryService).Upsert(req.ID, req.Hostname, req.Address)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, req)
		})
	}
}

func NewConsoleController(serverRepo proxyService.ServerRepositoryService) *ConsoleController {
	return &ConsoleController{serverRepo: serverRepo}
}
