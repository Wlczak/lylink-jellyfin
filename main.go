package main

import (
	"flag"
	"fmt"
	"net/http"

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

	headless := flag.Bool("headless", true, "Run in headless mode without desktop GUI.")
	flag.Parse()

	if *headless {
		api.RunHttpServer(srv)
		return
	}

	go api.RunHttpServer(srv)
	runApp(r, srv)
}
