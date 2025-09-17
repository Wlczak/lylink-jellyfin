package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func NewApi(username string, password string) (*Api, error) {

	request_body, err := json.Marshal(map[string]string{
		"username": username,
		"Pw":       password,
	})

	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("POST", "http://localhost:8096/Users/AuthenticateByName", bytes.NewReader(request_body))

	request.Header.Add("Content-Type", "application/json")

	connectionName := "lylink_jellyfin"
	request.Header.Add("X-Emby-Authorization", "Emby UserId=\""+username+"\", Client=\""+connectionName+"\", Device=\""+connectionName+"\", DeviceId=\""+connectionName+"\", Version=\"1.0\", Token=\"\"")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("request failed")
	}

	var body []byte

	body, _ = io.ReadAll(response.Body)

	fmt.Println(string(body))

	err = response.Body.Close()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("request closure failed")
	}

	return &Api{}, nil
}
