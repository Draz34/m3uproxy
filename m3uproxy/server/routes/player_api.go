package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Draz34/m3uproxy/config"
	"github.com/tidwall/sjson"
)

func PlayerApiRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/player_api.php", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}
		Action := r.FormValue("action")
		serieNum := r.FormValue("series_id")
		vodNum := r.FormValue("vod_id")

		formData := url.Values{
			"username":  {config.Xtream.Username},
			"password":  {config.Xtream.Password},
			"action":    {Action},
			"series_id": {serieNum},
			"vod_id":    {vodNum},
		}

		urlString := "http://" + config.Xtream.Hostname + ":" + strconv.Itoa(int(config.Xtream.Port)) + "/player_api.php"

		resp, err := http.PostForm(urlString, formData)
		if err != nil {
			print(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			print(err)
		}

		bodyStr := string(body)

		//fmt.Println(string(body))

		urlRequest := urlString
		if Action != "" {
			urlRequest = urlString + "?action=" + Action + "&series_id=" + serieNum + "&vod_id=" + vodNum
		} else {
			//fix json errors
			var re = regexp.MustCompile(`"(.*)": ([^"].*),`)
			bodyStr = re.ReplaceAllString(bodyStr, `"$1": "$2",`)
			bodyStr = strings.Replace(bodyStr, `90"LAR`, `90 LAR`, -1)

			//Modification de la réponse
			Username := r.FormValue("username")
			Password := r.FormValue("password")

			//log.Printf("Username = %s\n", Username)
			//log.Printf("Password = %s\n", Password)

			/*
				config.Xtream.Hostname = gjson.Get(bodyStr, "server_info.url").String()
				port, _ := strconv.Atoi(gjson.Get(bodyStr, "server_info.port").String())
				config.Xtream.Port = uint16(port)
			*/

			bodyStr, _ = sjson.Set(bodyStr, "user_info.username", Username)
			bodyStr, _ = sjson.Set(bodyStr, "user_info.password", Password)

			bodyStr, _ = sjson.Set(bodyStr, "server_info.url", config.Server.Hostname)
			bodyStr, _ = sjson.Set(bodyStr, "server_info.port", strconv.Itoa(int(config.Server.Port)))
			bodyStr, _ = sjson.Delete(bodyStr, "server_info.https_port")
			bodyStr, _ = sjson.Set(bodyStr, "server_info.server_protocol", "http")

			//Fin modification de la réponse
		}

		fmt.Println(urlRequest)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(bodyStr))
	}
}
