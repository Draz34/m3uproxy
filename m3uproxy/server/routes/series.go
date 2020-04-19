package routes

import (
	"net/http"
	"strconv"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/gorilla/mux"
)

func SeriesRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/series/{username}/{password}/{id}.{ext}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		channelNumber := vars["id"]
		ext := vars["ext"]
		username := vars["username"]
		password := vars["password"]

		if db.GetUser(username, password).ID <= 0 {
			w.WriteHeader(401)
			return
		}

		http.Redirect(w, r, "http://"+config.Server.Hostname+":"+strconv.Itoa(int(config.Server.Port))+"/channels/"+username+"/"+password+"/"+channelNumber+"?location=http%3A%2F%2F"+config.Xtream.Hostname+"%3A"+strconv.Itoa(int(config.Xtream.Port))+"%2Fseries%2F"+config.Xtream.Username+"%2F"+config.Xtream.Password+"%2F"+channelNumber+"."+ext, 301)
	}
}
