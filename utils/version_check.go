package utils

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
)

func HasUpdate() (bool, string, error) {
	resp, err := http.Get("https://api.github.com/repos/wlczak/lylink-jellyfin/releases/latest")
	if err != nil {
		return false, "", err
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return false, "", err
	}

	err = resp.Body.Close()

	if err != nil {
		return false, "", err
	}

	bi, ok := debug.ReadBuildInfo()

	return release.TagName != bi.Main.Version, release.TagName, nil
}
