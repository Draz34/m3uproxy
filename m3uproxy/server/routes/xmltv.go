package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
)

func XmltvRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/xmltv.php", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		if db.GetUser(username, password).ID <= 0 {
			w.WriteHeader(401)
			return
		}

		formData := url.Values{
			"username": {config.Xtream.Username},
			"password": {config.Xtream.Password},
		}

		urlString := "http://" + config.Xtream.Hostname + ":" + strconv.Itoa(int(config.Xtream.Port)) + "/xmltv.php"

		resp, err := http.PostForm(urlString, formData)
		if err != nil {
			print(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			print(err)
		}

		//fmt.Println(string(body))

		//Modification de la réponse
		bodyStr := string(body)
		bodyStr = strings.Replace(bodyStr, config.Xtream.Hostname+":"+strconv.Itoa(int(config.Xtream.Port)), config.Server.Hostname+":"+strconv.Itoa(int(config.Server.Port)), -1)

		fmt.Println(urlString)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(bodyStr))
	}
}
