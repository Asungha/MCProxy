package webui

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	metricService "mc_reverse_proxy/src/metric/service"
	proxyService "mc_reverse_proxy/src/proxy/service"
	"mc_reverse_proxy/src/webui/backend/common"
	"net/http"
	"strings"

	"embed"

	gin "github.com/gin-gonic/gin"
)

//go:embed frontend/build/*
var staticFiles embed.FS

//go:embed frontend/build/_next/static/chunks/*.*
var chunksFiles embed.FS

//go:embed frontend/build/_next/static/css/*.*
var cssFiles embed.FS

type WebUI struct {
	backend common.HTTPController
}

func NewWebUI(metricCollector *metricService.MetricService, serverRepo proxyService.ServerRepositoryService) *WebUI {
	return &WebUI{backend: NewHTTPBackend(metricCollector, serverRepo)}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// type StaticMapper struct {
// 	srcpath string
// 	router  *gin.Engine
// }

// func (m *StaticMapper) Map(route string) {
// 	r := route
// 	if route == "/" {
// 		abslutePath := m.srcpath + "/index.html"
// 		log.Printf("------------%s %s", route, abslutePath)
// 		abslutePath = strings.ReplaceAll(abslutePath, "//", "/")
// 		m.router.NoRoute(func(ctx *gin.Context) {
// 			ctx.FileFromFS(abslutePath, http.FS(htmlFiles))
// 		})
// 		return
// 	}
// 	abslutePath := m.srcpath + r + ".html"
// 	log.Printf("%s %s", route, abslutePath)
// 	abslutePath = strings.ReplaceAll(abslutePath, "//", "/")
// 	// m.router.LoadHTMLGlob(abslutePath)
// 	// m.router.StaticFile(route, abslutePath)
// 	m.router.GET(route, func(ctx *gin.Context) {
// 		ctx.FileFromFS(abslutePath, http.FS(htmlFiles))
// 	})
// }

// func NewStaticMapper(srcpath string, router *gin.Engine) *StaticMapper {
// 	return &StaticMapper{srcpath: srcpath, router: router}
// }

// func rewriteFilePath(filePath string) string {
// 	// Example: replace "oldpath" with "newpath"
// 	if filePath == "/oldpath" {
// 		filePath = "/newpath"
// 	}
// 	return filePath
// }

func (w *WebUI) Serve(listenAddr string) {
	// rootDir := "frontend/build"
	serverRoot, err := fs.Sub(staticFiles, "frontend/build")
	if err != nil {
		log.Fatal(err)
	}
	frouter := gin.Default()
	frouter.Use(CORSMiddleware())
	frouter.Use(func(ctx *gin.Context) {
		oldPath := ctx.Request.URL.Path
		if oldPath == "/" {
			ctx.Next()
			return
		} else if frags := strings.Split(oldPath, "."); len(frags) == 1 {
			ctx.Request.URL.Path = fmt.Sprintf("%s.html", oldPath)
			ctx.Redirect(http.StatusMovedPermanently, ctx.Request.URL.Path)
			return
		}
		ctx.Next()
	})
	// router.LoadHTMLGlob("webui/frontend/build/index.html")
	// router.LoadHTMLGlob("webui/frontend/build/dashboard.html")

	// router.Static("/_next/static", "webui/frontend/build/_next/static")
	frouter.StaticFS("/", http.FS(serverRoot))
	// chunks, err := fs.Sub(chunksFiles, "frontend/build/_next/static/chunks")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// frouter.StaticFS("/_next/static/chunks", http.FS(chunks))

	// css, err := fs.Sub(cssFiles, "frontend/build/_next/static/css")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// frouter.StaticFS("/_next/static/css", http.FS(css))
	// router.StaticFile("/favicon.ico", "webui/frontend/build/favicon.ico")

	// export const paths = {
	// 	home: '/',
	// 	auth: { signIn: '/auth/sign-in', signUp: '/auth/sign-up', resetPassword: '/auth/reset-password' },
	// 	dashboard: {
	// 	  overview: '/dashboard',
	// 	  server: '/dashboard/server',
	// 	  account: '/dashboard/account',
	// 	  customers: '/dashboard/customers',
	// 	  integrations: '/dashboard/integrations',
	// 	  settings: '/dashboard/settings',
	// 	},
	// 	errors: { notFound: '/errors/not-found' },
	//   } as const;

	// mapper := NewStaticMapper("frontend/build/", router)

	// mapper.Map("/dashboard")
	// mapper.Map("/auth/sign-in")
	// mapper.Map("/auth/sign-up")
	// mapper.Map("/auth/reset-password")
	// mapper.Map("/dashboard")
	// mapper.Map("/dashboard/server")
	// mapper.Map("/errors/not-found")

	brouter := gin.Default()
	brouter.Use(CORSMiddleware())
	w.backend.Config(brouter)

	ctx, cancle := context.WithCancel(context.Background())

	go func() {
		frouter.Run("0.0.0.0:8082")
		cancle()
	}()
	go func() {
		brouter.Run("0.0.0.0:8080")
		cancle()
	}()

	<-ctx.Done()
}
