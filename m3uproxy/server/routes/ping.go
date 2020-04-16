package routes

import (
	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/server/webutils"
	"net/http"
)

var bytes = []byte("pong")

func PingRouter(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {
	return "/ping", func(w http.ResponseWriter, r *http.Request) {
		webutils.Success(bytes, w)
	}
}
