package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Wlczak/lylink-jellyfin/api"
	"github.com/Wlczak/lylink-jellyfin/desktop"
	"github.com/Wlczak/lylink-jellyfin/logs"
	"github.com/gin-gonic/gin"

	_ "embed"
)

type GetTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetPlaybackInfoRequest struct {
	AccessToken string `json:"token"`
}

type GetMediaInfoRequest struct {
	AccessToken string `json:"token"`
}

//go:embed Icon.png
var icon []byte

func runApp() {
	a := desktop.Init(icon)

	a.Run()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	setupRoutes(r)
	go runHttpServer(r)

	runApp()
}

func setupRoutes(r *gin.Engine) {
	zap := logs.GetLogger()

	r.GET("/handshake", func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, "hand shaken")
	})

	r.POST("/getToken", func(c *gin.Context) {
		bodyReader := c.Request.Body

		body, err := io.ReadAll(bodyReader)
		if err != nil {
			zap.Error(err.Error())
		}

		err = bodyReader.Close()

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
			if err.Error() == "no media playing" {
				c.JSON(http.StatusOK, nil)
				return
			} else {
				zap.Error(err.Error())
			}
		}

		c.JSON(http.StatusOK, sessions)
	})

	r.POST("/Item/:id", func(c *gin.Context) {
		bodyReader := c.Request.Body
		body, err := io.ReadAll(bodyReader)
		mediaId := c.Param("id")

		if err != nil {
			zap.Error(err.Error())
		}

		err = bodyReader.Close()

		if err != nil {
			zap.Error(err.Error())
		}

		r := GetMediaInfoRequest{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			zap.Error(err.Error())
		}

		apiObj := api.NewApi(r.AccessToken)

		episodeInfo, err := apiObj.GetEpisodeInfo(mediaId)

		if err != nil {
			zap.Error(err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		var response api.GetMediaInfoResponse

		seasonInfo, err := apiObj.GetSeasonInfo(episodeInfo.ParentId)
		response = api.GetMediaInfoResponse{Id: episodeInfo.Id, Name: episodeInfo.Name, Type: episodeInfo.Type, SeriesName: episodeInfo.SeriesName, IndexNumber: episodeInfo.IndexNumber, ParentIndexNumber: episodeInfo.ParentIndexNumber}

		if err != nil {
			zap.Error(err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		response.SeasonId = seasonInfo.Id
		seriesInfo, err := apiObj.GetSeriesInfo(seasonInfo.ParentId)

		if err != nil {
			zap.Error(err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		response.SeriesId = seriesInfo.Id

		c.JSON(http.StatusOK, response)
	})
}

func runHttpServer(r *gin.Engine) {
	zap := logs.GetLogger()

	fmt.Println("Listening on port :8040")

	err := r.Run(":8040")

	if err != nil {
		zap.Error(err.Error())
	}
}
