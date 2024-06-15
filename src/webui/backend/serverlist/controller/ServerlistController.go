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

type AddServerRecord struct {
	Hostname string `json:hostname`
	Address  string `json:address`
}

type DeleteServerRecord struct {
	ID int `json:id`
}

type ServerlistController struct {
	serverRepo proxyService.ServerRepositoryService
}

func (c *ServerlistController) Config(router *gin.Engine) {
	r := router.Group("/console")
	{
		r.GET("/server-list", func(ctx *gin.Context) {
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
		r.POST("/server-list", func(ctx *gin.Context) {
			var req AddServerRecord
			if err := ctx.BindJSON(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Printf("%v", req)
			err := c.serverRepo.(proxyService.UpdatableRepositoryService).Insert(req.Hostname, req.Address)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, req)
		})
		r.PUT("/server-list", func(ctx *gin.Context) {
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
		r.DELETE("/server-list", func(ctx *gin.Context) {
			var req DeleteServerRecord
			if err := ctx.BindJSON(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Printf("%v", req)
			err := c.serverRepo.(proxyService.UpdatableRepositoryService).Delete(req.ID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, req)
		})
	}
}

func NewServerlistController(serverRepo proxyService.ServerRepositoryService) *ServerlistController {
	return &ServerlistController{serverRepo: serverRepo}
}