package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/Draz34/m3uproxy/server/webutils"
	"github.com/gorilla/mux"
)

func MovieRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/movie/{username}/{password}/{id}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		channelNumber := vars["id"]
		username := vars["username"]
		password := vars["password"]

		if db.GetUser(username, password).ID <= 0 {
			w.WriteHeader(401)
			return
		}

		channel, err := db.LookupChannel(channelNumber)
		trackRedirects := false
		var urlIptv string
		if err != nil {
			urlIptv = "http://" + config.Xtream.Hostname + ":" + strconv.Itoa(int(config.Xtream.Port)) + "/movie/" + config.Xtream.Username + "/" + config.Xtream.Password + "/" + channelNumber
			log.Printf("Register Channel for %s", urlIptv)
			channel, _ = db.RegisterChannel(urlIptv)
			//log.Printf("%+v\n", channel)

			//Si l'url n'est pas en mémoire on trace les redirections
			trackRedirects = true
		}

		redirectUrl := "http://" + config.Server.Hostname + ":" + strconv.Itoa(int(config.Server.Port)) + "/channels/" + username + "/" + password + "/" + channel.Id
		log.Printf("Redirect to %s", redirectUrl)

		if trackRedirects {
			webutils.TracingRedirect(urlIptv)
		}

		http.Redirect(w, r, redirectUrl, 302)
	}
}
