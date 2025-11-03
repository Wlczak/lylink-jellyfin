package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Wlczak/lylink-jellyfin/api"
	"github.com/Wlczak/lylink-jellyfin/config"
	"github.com/Wlczak/lylink-jellyfin/desktop"
	"github.com/gin-gonic/gin"

	_ "embed"
)

//go:embed Icon.png
var icon []byte

func runApp(r *gin.Engine, srv *http.Server) {
	a := desktop.Init(icon, r, srv)

	a.Run()
}

func main() {
	conf := config.GetConfig()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: r,
	}

	api.SetupRoutes(r)

	for _, arg := range os.Args {
		if arg == "--headless" {
			api.RunHttpServer(srv)
			return
		}
	}
	go api.RunHttpServer(srv)
	runApp(r, srv)
}
