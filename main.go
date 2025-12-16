package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Wlczak/lylink-jellyfin/api"
	"github.com/Wlczak/lylink-jellyfin/config"
	"github.com/Wlczak/lylink-jellyfin/desktop"
	"github.com/Wlczak/lylink-jellyfin/utils"
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
	headless := flag.Bool("headless", false, "Run in headless mode without desktop GUI.")
	versionCheck := flag.Bool("version", false, "Check for updates and exit.")
	versionCheckShort := flag.Bool("v", false, "Check for updates and exit.")
	flag.Parse()
	if versionCheck != nil && *versionCheck || versionCheckShort != nil && *versionCheckShort {
		hasUpdate, versionName, err := utils.HasUpdate()

		if err != nil {
			fmt.Println("Error checking for updates")
			return
		}

		if hasUpdate {
			fmt.Println("New version " + versionName + " available")
		} else {
			fmt.Println("No updates available")
		}
		return
	}

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

	if *headless {
		api.RunHttpServer(srv)
		return
	}

	go api.RunHttpServer(srv)
	runApp(r, srv)
}
