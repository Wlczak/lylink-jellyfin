package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

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

	request_body, err := json.Marshal(map[string]string{
		"Username": username,
		"Pw":       password,
	})

	if err != nil {
		return nil, err
	}

	request := newRequest(http.MethodPost, "http://localhost:8096/Users/AuthenticateByName", username, bytes.NewBuffer(request_body))

	body, _, err := execRequest(request)

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

	request := newRequest(http.MethodGet, "http://localhost:8096/Sessions", "", nil)

	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, _, err := execRequest(request)

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

func (api *Api) GetMediaInfo(mediaSourceId string) (MediaInfo, error) {
	zap := logs.GetLogger()

	request := newRequest(http.MethodGet, "http://localhost:8096/Item/"+mediaSourceId, "", nil)

	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, _, err := execRequest(request)

	if err != nil {
		zap.Error(err.Error())
		return MediaInfo{}, err
	}

	var mediaInfo MediaInfo
	err = json.Unmarshal(body, &mediaInfo)

	if err != nil {
		return MediaInfo{}, err
	}
	return mediaInfo, nil
}

func NewApi(token string) *Api {
	return &Api{AccessToken: token, Username: "guest"}
}
