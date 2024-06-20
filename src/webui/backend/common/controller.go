package common

import "github.com/gin-gonic/gin"

type HTTPController interface {
	Config(*gin.Engine)
}
