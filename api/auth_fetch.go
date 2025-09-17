package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func NewApi(username string, password string) (*Api, error) {

	request_body := []byte("{\"username\":\"" + username + "\",\"Pw\":\"" + password + "\"}")

	request, _ := http.NewRequest("POST", "http://localhost:8096/Users/AuthenticateByName", bytes.NewReader(request_body))

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Emby-Authorization", "Emby UserId=\"Hieroglyph.Admin\", Client=\"media_cleaner\", Device=\"media_cleaner\", DeviceId=\"media_cleaner\", Version=\"0.5\", Token=\"\"")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("request failed")
	}
	defer response.Body.Close()

	var body []byte

	body, _ = io.ReadAll(response.Body)

	fmt.Println(string(body))

	return &Api{}, nil
}
