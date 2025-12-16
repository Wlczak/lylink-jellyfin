package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime/debug"
	"strings"
)

// HasUpdate checks if a new version is available
//
// Return values are: bolean (true if update is available), string newVersionName, string currentVersionName, error err
func HasUpdate() (bool, string, string, error) {
	resp, err := http.Get("https://api.github.com/repos/wlczak/lylink-jellyfin/releases/latest")
	if err != nil {
		return false, "", "", err
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return false, "", "", err
	}

	err = resp.Body.Close()

	if err != nil {
		return false, "", "", err
	}

	bi, ok := debug.ReadBuildInfo()

	if !ok {
		return false, "", "", errors.New("failed to read build info")
	}

	return release.TagName != strings.Split(bi.Main.Version, "+")[0], release.TagName, strings.Split(bi.Main.Version, "+")[0], nil
}

func GetCurrentVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return strings.Split(bi.Main.Version, "+")[0]
}
