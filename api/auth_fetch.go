package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func newRequest(method string, url string, username string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(method, url, body)

	if method == http.MethodPost {
		request.Header.Add("Content-Type", "application/json")
	}

	connectionName := "lylink_jellyfin"
	if username != "" {
		request.Header.Add("X-Emby-Authorization", "Emby UserId=\""+username+"\", Client=\""+connectionName+"\", Device=\""+connectionName+"\", DeviceId=\""+connectionName+"\", Version=\"1.0\", Token=\"\"")
	}

	return request
}

func execRequest(request *http.Request) (body []byte, response *http.Response, err error) {
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, nil, errors.New("request failed")
	}

	body, _ = io.ReadAll(response.Body)
	return body, response, nil
}

func NewApi(username string, password string) (*Api, error) {

	request_body, err := json.Marshal(map[string]string{
		"Username": username,
		"Pw":       password,
	})

	if err != nil {
		return nil, err
	}

	request := newRequest(http.MethodPost, "http://localhost:8096/Users/AuthenticateByName", username, bytes.NewBuffer(request_body))

	body, _, err := execRequest(request)

	//fmt.Println(string(body))

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("request closure failed")
	}

	var authResponse AuthResponse

	err = json.Unmarshal(body, &authResponse)

	return &Api{Username: username, AccessToken: authResponse.AccessToken}, nil
}

func (api *Api) GetPlaybackInfo() {
	request := newRequest(http.MethodGet, "http://localhost:8096/Sessions", "", nil)
	fmt.Println(api.AccessToken)
	request.Header.Set("Authorization", "MediaBrowser Token="+api.AccessToken)

	body, response, _ := execRequest(request)

	fmt.Println(response.StatusCode)
	fmt.Println(len(body), string(body))
}
