package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/Wlczak/lylink-jellyfin/api"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api, err := api.NewApi("username", "password")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		item := api.GetPlaybackInfo()

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(fmt.Sprint(item.PlayState.PositionTicks)))
	})

	list, _ := net.Listen("tcp", ":8040")
	fmt.Println("Listening on port 8040")
	http.Serve(list, nil)
}
