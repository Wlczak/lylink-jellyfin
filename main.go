package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Wlczak/lylink-jellyfin/api"
	"github.com/Wlczak/lylink-jellyfin/logs"
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
	zap := logs.GetLogger()
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

		body, err := io.ReadAll(bodyReader)
		err = bodyReader.Close()
		if err != nil {
			zap.Error(err.Error())
		}

		if err != nil {
			zap.Error(err.Error())
		}

		u := GetTokenRequest{}
		err = json.Unmarshal(body, &u)
		if err != nil {
			zap.Error(err.Error())
		}

		api, err := api.GetToken(u.Username, u.Password)

		if err != nil {
			zap.Error(err.Error())
		}

		// c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, api.AccessToken)
	})

	r.POST("/getPlaybackInfo", func(c *gin.Context) {
		bodyReader := c.Request.Body
		body, err := io.ReadAll(bodyReader)

		if err != nil {
			zap.Error(err.Error())
		}

		err = bodyReader.Close()

		if err != nil {
			zap.Error(err.Error())
		}

		r := GetPlaybackInfoRequest{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			zap.Error(err.Error())
		}

		api := api.NewApi(r.AccessToken)

		sessions, err := api.GetPlaybackInfo()

		if err != nil {
			zap.Error(err.Error())
		}

		c.JSON(http.StatusOK, sessions)
	})

	fmt.Println("Listening on port :8040")

	err := r.Run(":8040")

	if err != nil {
		zap.Error(err.Error())
	}
}
