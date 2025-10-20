package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Wlczak/lylink-jellyfin/config"
	"github.com/Wlczak/lylink-jellyfin/logs"
)

func newRequest(method string, url string, username string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(method, url, body)

	if method == http.MethodPost {
		request.Header.Add("Content-Type", "application/json")
	}

	connectionName := "lylink_jellyfin"
	if username != "" {
		// Warning this is a very magic string, any alteration may break the auth process.
		// You have beeen warned proceed with caution
		request.Header.Add("X-Emby-Authorization", "Emby Client=\""+connectionName+"\", Device=\""+connectionName+"\", DeviceId=\""+connectionName+"\", Version=\"1.0\"")
	}

	return request
}

func execRequest(request *http.Request) (body []byte, response *http.Response, err error) {
	zap := logs.GetLogger()
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		zap.Error(err.Error())
		return nil, nil, errors.New("request failed")
	}

	body, err = io.ReadAll(response.Body)
	if err != nil {
		zap.Error(err.Error())
		return nil, nil, errors.New("request closure failed")
	}
	return body, response, nil
}

func GetToken(username string, password string) (*Api, error) {
	zap := logs.GetLogger()
	conf := config.GetConfig()

	request_body, err := json.Marshal(map[string]string{
		"Username": username,
		"Pw":       password,
	})

	if err != nil {
		return nil, err
	}

	request := newRequest(http.MethodPost, conf.JellyfinServerUrl+"/Users/AuthenticateByName", username, bytes.NewBuffer(request_body))

	body, response, err := execRequest(request)

	if response.StatusCode != 200 {
		zap.Error(err.Error())
		return nil, errors.New("auth failed")
	}

	if err != nil {
		zap.Error(err.Error())
		return nil, errors.New("request closure failed")
	}

	var authResponse AuthResponse

	err = json.Unmarshal(body, &authResponse)

	if err != nil {
		zap.Error(err.Error())
	}

	if authResponse.AccessToken == "" {
		return nil, errors.New("auth failed")
	}

	return &Api{Username: username, AccessToken: authResponse.AccessToken}, nil
}

func (api *Api) GetPlaybackInfo() ([]SessionItem, error) {
	zap := logs.GetLogger()
	conf := config.GetConfig()

	request := newRequest(http.MethodGet, conf.JellyfinServerUrl+"/Sessions", "", nil)

	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, response, err := execRequest(request)

	if response.StatusCode != 200 {
		err = errors.New(response.Status)
		zap.Error(err.Error())
		return nil, err
	}

	if err != nil {
		zap.Error(err.Error())
		return nil, err
	}

	var items []SessionItem
	err = json.Unmarshal(body, &items)

	if err != nil {
		return nil, err
	}

	var activeItems []SessionItem
	var hasMediaPlaying = false
	for _, item := range items {
		if item.PlayState.MediaSourceId != "" {
			activeItems = append(activeItems, item)
			hasMediaPlaying = true
		}
	}

	if hasMediaPlaying {
		return activeItems, nil

	}
	return nil, errors.New("no media playing")
}

func (api *Api) GetEpisodeInfo(mediaSourceId string) (EpisodeInfo, error) {
	zap := logs.GetLogger()
	conf := config.GetConfig()

	request := newRequest(http.MethodGet, conf.JellyfinServerUrl+"/Items/"+mediaSourceId, "", nil)

	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, response, err := execRequest(request)

	if response.StatusCode != 200 {
		err = errors.New(response.Status)
		zap.Error(err.Error())
		return EpisodeInfo{}, err
	}

	if err != nil {
		zap.Error(err.Error())
		return EpisodeInfo{}, err
	}

	var mediaInfo MediaInfo
	err = json.Unmarshal(body, &mediaInfo)

	if err != nil {
		return EpisodeInfo{}, err
	}

	if mediaInfo.Type == "" {
		return EpisodeInfo{}, errors.New("no media info")
	}

	var info EpisodeInfo

	switch mediaInfo.Type {
	case "Episode":
		info = EpisodeInfo{}
		err = json.Unmarshal(body, &info)

	default:
		err = errors.New("media not episode type")
		zap.Error(err.Error())
		return EpisodeInfo{}, err
	}

	if err != nil {
		return EpisodeInfo{}, err
	}

	return info, nil
}

func (api *Api) GetSeasonInfo(mediaSourceId string) (SeasonInfo, error) {
	zap := logs.GetLogger()
	conf := config.GetConfig()

	request := newRequest(http.MethodGet, conf.JellyfinServerUrl+"/Items/"+mediaSourceId, "", nil)

	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, response, err := execRequest(request)

	if response.StatusCode != 200 {
		err = errors.New(response.Status)
		zap.Error(err.Error())
		return SeasonInfo{}, err
	}

	if err != nil {
		zap.Error(err.Error())
		return SeasonInfo{}, err
	}

	var mediaInfo MediaInfo
	err = json.Unmarshal(body, &mediaInfo)

	if err != nil {
		return SeasonInfo{}, err
	}

	if mediaInfo.Type == "" {
		return SeasonInfo{}, errors.New("no media info")
	}

	var info SeasonInfo

	switch mediaInfo.Type {
	case "Season":
		info = SeasonInfo{}
		err = json.Unmarshal(body, &info)

	default:
		err = errors.New("media not Season type")
		zap.Error(err.Error())
		return SeasonInfo{}, err
	}

	if err != nil {
		return SeasonInfo{}, err
	}

	return info, nil
}

func (api *Api) GetSeriesInfo(mediaSourceId string) (SeriesInfo, error) {
	zap := logs.GetLogger()
	conf := config.GetConfig()

	request := newRequest(http.MethodGet, conf.JellyfinServerUrl+"/Items/"+mediaSourceId, "", nil)

	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, response, err := execRequest(request)

	if response.StatusCode != 200 {
		err = errors.New(response.Status)
		zap.Error(err.Error())
		return SeriesInfo{}, err
	}

	if err != nil {
		zap.Error(err.Error())
		return SeriesInfo{}, err
	}

	var mediaInfo MediaInfo
	err = json.Unmarshal(body, &mediaInfo)

	if err != nil {
		return SeriesInfo{}, err
	}

	if mediaInfo.Type == "" {
		return SeriesInfo{}, errors.New("no media info")
	}

	var info SeriesInfo

	switch mediaInfo.Type {
	case "Series":
		info = SeriesInfo{}
		err = json.Unmarshal(body, &info)

	default:
		err = errors.New("media not Series type")
		zap.Error(err.Error())
		return SeriesInfo{}, err
	}

	if err != nil {
		return SeriesInfo{}, err
	}

	return info, nil
}

func (api *Api) GetEpisodeList(seriesId string) ([]EpisodeInfo, error) {
	zap := logs.GetLogger()
	conf := config.GetConfig()

	request := newRequest(http.MethodGet, conf.JellyfinServerUrl+"/Shows/"+seriesId+"/Episodes", "", nil)

	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, response, err := execRequest(request)

	if response.StatusCode != 200 {
		err = errors.New(response.Status)
		zap.Error(err.Error())
		return []EpisodeInfo{}, err
	}

	if err != nil {
		zap.Error(err.Error())
		return []EpisodeInfo{}, err
	}

	var episodeList EpisodeList
	err = json.Unmarshal(body, &episodeList)

	if err != nil {
		return []EpisodeInfo{}, err
	}

	if episodeList.Items == nil {
		return []EpisodeInfo{}, errors.New("no media info")
	}

	return episodeList.Items, nil
}

func NewApi(token string) *Api {
	return &Api{AccessToken: token, Username: "guest"}
}
