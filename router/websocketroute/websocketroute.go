package websocketroute

import (
	"TestChat1/controller/websocketcontroller"
	"net/http"
)

type websocketFunc func()

func StartWebSocketRoute() {
	http.HandleFunc("/ws", websocketcontroller.FirstPage)
	http.ListenAndServe(":8087", nil)
}
