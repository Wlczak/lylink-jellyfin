package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Wlczak/lylink-jellyfin/logs"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
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

		api, err := GetToken(u.Username, u.Password)

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

		api := NewApi(r.AccessToken)

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

	r.POST("/Episode/WithParents/:id", func(c *gin.Context) {
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

		apiObj := NewApi(r.AccessToken)

		episodeInfo, err := apiObj.GetEpisodeInfo(mediaId)

		if err != nil {
			zap.Error(err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		var response GetMediaInfoResponse

		seasonInfo, err := apiObj.GetSeasonInfo(episodeInfo.ParentId)
		response = GetMediaInfoResponse{Id: episodeInfo.Id, Name: episodeInfo.Name, Type: episodeInfo.Type, SeriesName: episodeInfo.SeriesName, IndexNumber: episodeInfo.IndexNumber, ParentIndexNumber: episodeInfo.ParentIndexNumber}

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

func RunHttpServer(r *http.Server) {
	zap := logs.GetLogger()

	fmt.Println("Listening on port :8040")

	err := r.ListenAndServe()

	if err != nil {
		zap.Error(err.Error())
		if err.Error() == "http: Server closed" {
			return
		} else {
			panic(err)
		}

	}
}
