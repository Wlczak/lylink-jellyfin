package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/handshake", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, "hand shaken")
	})

	fmt.Println("Listening on port :8040")
	r.Run(":8040")
}
