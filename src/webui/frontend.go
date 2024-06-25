package webui

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"mc_reverse_proxy/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/build/*
var staticFiles embed.FS

type HTTPFrontend struct {
	address         string
	backendHostname string
}

func (f *HTTPFrontend) Config(engine *gin.Engine) {
	serverRoot, err := fs.Sub(staticFiles, "frontend/build")
	if err != nil {
		log.Fatal(err)
	}
	engine.Use(CORSMiddleware())
	engine.Use(func(ctx *gin.Context) {
		oldPath := ctx.Request.URL.Path
		if oldPath == "/api.json" {
			ctx.JSON(http.StatusOK, gin.H{"api_address": f.backendHostname})
			ctx.Done()
			return
		} else if oldPath == "/" {
			ctx.Next()
			return
		} else if frags := strings.Split(oldPath, "."); len(frags) == 1 {
			ctx.Request.URL.Path = fmt.Sprintf("%s.html", oldPath)
			ctx.Redirect(http.StatusMovedPermanently, ctx.Request.URL.Path)
			return
		}
		ctx.Next()
	})
	engine.StaticFS("/", http.FS(serverRoot))

}

func (f *HTTPFrontend) Serve() error {
	engine := gin.Default()
	f.Config(engine)
	utils.FLog.HTTPFrontend("Start server on %s", f.address)
	return engine.Run(f.address)
}

func NewHTTPFrontend(address string, backendHostname string) *HTTPFrontend {
	gin.SetMode(gin.ReleaseMode)
	f := &HTTPFrontend{
		address:         address,
		backendHostname: backendHostname,
	}
	return f
}
