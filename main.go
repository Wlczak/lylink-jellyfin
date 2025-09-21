package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Wlczak/lylink-jellyfin/api"
	"github.com/gin-gonic/gin"
)

type GetTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/handshake", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, "hand shaken")
	})

	r.POST("/getToken", func(c *gin.Context) {
		bodyReader := c.Request.Body
		defer bodyReader.Close()
		body, _ := io.ReadAll(bodyReader)

		u := GetTokenRequest{}
		json.Unmarshal(body, &u)

		api, _ := api.NewApi(u.Username, u.Password)

		c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, api.AccessToken)
	})

	fmt.Println("Listening on port :8040")
	r.Run(":8040")
}
