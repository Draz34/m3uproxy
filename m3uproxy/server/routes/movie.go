package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/gorilla/mux"
)

func MovieRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/movie/{username}/{password}/{id}.{ext}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		channelNumber := vars["id"]
		ext := vars["ext"]
		username := vars["username"]
		password := vars["password"]

		if db.GetUser(username, password).ID <= 0 {
			w.WriteHeader(401)
			return
		}

		channel, err := db.LookupChannel(channelNumber)
		if err != nil {
			urlIptv := "http://" + config.Xtream.Hostname + ":" + strconv.Itoa(int(config.Xtream.Port)) + "/movie/" + config.Xtream.Username + "/" + config.Xtream.Password + "/" + channelNumber + "." + ext
			log.Printf("Register Channel for %s", urlIptv)
			channel, _ = db.RegisterChannel(urlIptv)
			log.Printf("%+v\n", channel)
		}

		redirectUrl := "http://" + config.Server.Hostname + ":" + strconv.Itoa(int(config.Server.Port)) + "/channels/" + username + "/" + password + "/" + channel.Id
		log.Printf("Redirect to %s", redirectUrl)

		http.Redirect(w, r, redirectUrl, 302)
	}
}
