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

		item, _ := api.GetPlaybackInfo()

		w.Header().Set("Content-Type", "text/plain")

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusOK)

		percentage := float64(item.PlayState.PositionTicks) / float64(item.NowPlayingItem.RunTimeTicks+1)
		w.Write([]byte(fmt.Sprintf("%f", percentage*100)))
	})

	list, _ := net.Listen("tcp", ":8040")
	fmt.Println("Listening on port 8040")
	http.Serve(list, nil)
}
