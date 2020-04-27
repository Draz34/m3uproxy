package routes

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/Draz34/m3uproxy/server/webutils"
	"github.com/gorilla/mux"
)

func LiveRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/live/{username}/{password}/{id}.{ext}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		channelNumber := vars["id"] + "." + vars["ext"]
		username := vars["username"]
		password := vars["password"]

		if db.GetUser(username, password).ID <= 0 {
			w.WriteHeader(401)
			return
		}

		channel, err := db.LookupChannel(channelNumber)

		var urlIptv string = "http://" + config.Xtream.Hostname + ":" + strconv.Itoa(int(config.Xtream.Port)) + "/live/" + config.Xtream.Username + "/" + config.Xtream.Password + "/" + channelNumber
		var urlIptvTs string = "http://" + config.Xtream.Hostname + ":" + strconv.Itoa(int(config.Xtream.Port)) + "/live/" + config.Xtream.Username + "/" + config.Xtream.Password + "/" + vars["id"] + ".ts"
		if err != nil {
			lastUrlIptv := webutils.TracingRedirect(urlIptvTs)

			if vars["ext"] == "m3u8" {
				log.Printf("urlIptv url now : %s", urlIptv)
				lastUrlIptv = strings.Replace(urlIptv, ".ts", ".m3u8", -1)
				urlIptv = lastUrlIptv
			}

			log.Printf("Register Channel for %s", urlIptv)
			channel, _ = db.RegisterChannel(urlIptv)
			//log.Printf("%+v\n", channel)
		} else {
			log.Printf("Update Channel for %s", urlIptv)
			channel, _ = db.UpdateChannel(channelNumber, urlIptv)
		}

		redirectUrl := "http://" + config.Server.Hostname + ":" + strconv.Itoa(int(config.Server.Port)) + "/channels/" + username + "/" + password + "/" + channel.Id
		log.Printf("Redirect to %s", redirectUrl)

		http.Redirect(w, r, redirectUrl, 302)
	}
}
