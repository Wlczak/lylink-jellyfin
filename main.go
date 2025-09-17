package main

import (
	"fmt"

	"github.com/Wlczak/lylink-jellyfin/api"
)

func main() {
	api, err := api.NewApi("username", "password")
	if err != nil {
		panic(err)
	}
	fmt.Println(api.AccessToken)
}
