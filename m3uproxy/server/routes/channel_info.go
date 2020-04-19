package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/Draz34/m3uproxy/server/webutils"
	"github.com/gorilla/mux"
)

func ChannelInfoRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {
	return "/channels/{username}/{password}/info/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["id"]
		username := vars["username"]
		password := vars["password"]

		if db.GetUser(username, password).ID <= 0 {
			w.WriteHeader(401)
			return
		}

		channel, err := db.LookupChannel(channelId)
		if err != nil {
			webutils.NotFound(w)
		}

		jsonChannel, err := json.Marshal(channel)
		if err != nil {
			webutils.InternalServerError("Error building jsonChannel from a db.Channel", err, w)
		}

		w.Header().Set("Content-Type", "application/json")
		webutils.Success(jsonChannel, w)
	}
}
