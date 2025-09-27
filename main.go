package main

import (
	"encoding/json"
	"errors"
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
	AccessToken string `json:"token"`
}

type GetMediaInfoRequest struct {
	AccessToken string `json:"token"`
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

		episodeInfo, err := apiObj.GetMediaInfo(mediaId)

		if err != nil {
			zap.Error(err.Error())
			c.Error(err)
			return
		}

		var seasonInfo api.Media
		var response api.GetMediaInfoResponse

		switch epInfo := episodeInfo.(type) {
		case api.EpisodeInfo:
			seasonInfo, err = apiObj.GetMediaInfo(epInfo.ParentId)
			response = api.GetMediaInfoResponse{Id: epInfo.Id, Name: epInfo.Name, Type: epInfo.Type, SeriesName: epInfo.SeriesName, IndexNumber: epInfo.IndexNumber, ParentIndexNumber: epInfo.ParentIndexNumber}
		default:
			err = errors.New("not correct type")
			zap.Error(err.Error())
			c.Error(err)
			return
		}

		var seriesInfo api.Media
		switch seasInfo := seasonInfo.(type) {
		case api.SeriesInfo:
			response.SeasonId = seasInfo.Id
			seriesInfo, err = apiObj.GetMediaInfo(seasInfo.Id)
		}

		switch serInfo := seriesInfo.(type) {
		case api.SeriesInfo:
			response.SeriesId = serInfo.Id
		}

		c.JSON(http.StatusOK, response)
	})

	fmt.Println("Listening on port :8040")

	err := r.Run(":8040")

	if err != nil {
		zap.Error(err.Error())
	}
}
