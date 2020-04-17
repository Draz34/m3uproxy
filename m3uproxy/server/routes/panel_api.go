package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Draz34/m3uproxy/config"
)

func PanelApiRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/panel_api.php", func(w http.ResponseWriter, r *http.Request) {

		formData := url.Values{
			"surname":  {config.Xtream.Surname},
			"username": {config.Xtream.Username},
			"password": {config.Xtream.Password},
		}

		urlString := "http://" + config.Xtream.Hostname + ":" + strconv.Itoa(int(config.Xtream.Port)) + "/panel_api.php"
		fmt.Println(urlString)

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

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
