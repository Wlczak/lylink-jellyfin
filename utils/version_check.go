package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

var Version = "dev"

func HasUpdate() (bool, string, error) {
	resp, err := http.Get("https://api.github.com/repos/wlczak/lylink-jellyfin/releases/latest")
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return false, "", err
	}
	bi, _ := debug.ReadBuildInfo()

	fmt.Println("Current version: ", bi.Main.Version)
	return release.TagName != Version, release.TagName, nil
}
