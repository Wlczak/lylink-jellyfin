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

type GetPlaybackInfoRequest struct {
	AccessToken string `json:"accessToken"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	r.GET("/handshake", func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, "hand shaken")
	})

	r.POST("/getToken", func(c *gin.Context) {
		bodyReader := c.Request.Body
		defer bodyReader.Close()
		body, _ := io.ReadAll(bodyReader)

		u := GetTokenRequest{}
		json.Unmarshal(body, &u)

		api, _ := api.GetToken(u.Username, u.Password)

		// c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, api.AccessToken)
	})

	r.POST("/getPlaybackInfo", func(c *gin.Context) {
		bodyReader := c.Request.Body
		defer bodyReader.Close()
		body, _ := io.ReadAll(bodyReader)

		r := GetPlaybackInfoRequest{}
		json.Unmarshal(body, &r)

		api := api.NewApi(r.AccessToken)

		sessions, _ := api.GetPlaybackInfo()

		c.JSON(http.StatusOK, sessions)
	})

	fmt.Println("Listening on port :8040")
	r.Run(":8040")
}
