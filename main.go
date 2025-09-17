package main

import "github.com/Wlczak/lylink-jellyfin/api"

func main() {
	_, err := api.NewApi("username", "password")
	if err != nil {
		panic(err)
	}
}
