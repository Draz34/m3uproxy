package routes

import (
	"net/http"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/server/webutils"
)

var resp = []byte("Welcome to m3u proxy")

func RootRouter(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {
	return "/", func(w http.ResponseWriter, r *http.Request) {
		webutils.Success(resp, w)
	}
}
