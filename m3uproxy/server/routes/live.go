package routes

import (
	"net/http"
	"strconv"

	"github.com/Draz34/m3uproxy/config"
	"github.com/gorilla/mux"
)

func LiveRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/live/{username}/{password}/{id}.ts", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		channelNumber := vars["id"]

		http.Redirect(w, r, "http://"+config.Server.Hostname+":"+strconv.Itoa(int(config.Server.Port))+"/channels/"+channelNumber+"?location=http%3A%2F%2F"+config.Xtream.Hostname+"%3A"+strconv.Itoa(int(config.Xtream.Port))+"%2Flive%2F"+config.Xtream.Username+"%2F"+config.Xtream.Password+"%2F"+channelNumber+".ts", 301)
	}
}