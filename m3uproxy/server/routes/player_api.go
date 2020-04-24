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
	"github.com/Draz34/m3uproxy/db"
	"github.com/tidwall/sjson"
)

func PlayerApiRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {

	return "/player_api.php", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}
		Action := r.FormValue("action")
		Username := r.FormValue("username")
		Password := r.FormValue("password")
		categorieNum := r.FormValue("category_id")
		streamNum := r.FormValue("stream_id")
		serieNum := r.FormValue("series_id")
		vodNum := r.FormValue("vod_id")
		Limit := r.FormValue("limit")

		bodyStr := `{
			"user_info": {
			  "auth": 0
			}
		  }`

		usr := db.GetUser(Username, Password)
		if usr.ID > 0 {
			formData := url.Values{
				"username":    {config.Xtream.Username},
				"password":    {config.Xtream.Password},
				"action":      {Action},
				"category_id": {categorieNum},
				"stream_id":   {streamNum},
				"series_id":   {serieNum},
				"vod_id":      {vodNum},
				"limit":       {Limit},
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

			//om := server.ordered.NewOrderedMap()
			//err = json.Unmarshal(body, om)

			//fmt.Print(om)

			bodyStr = string(body)

			//fmt.Println(string(body))

			urlRequest := urlString
			if Action != "" {
				urlRequest = urlString + "?action=" + Action + "&category_id=" + categorieNum + "&stream_id=" + streamNum + "&series_id=" + serieNum + "&vod_id=" + vodNum + "&limit=" + Limit
			} else {
				//fix json errors
				var re = regexp.MustCompile(`"(.*)": ([^"].*),`)
				bodyStr = re.ReplaceAllString(bodyStr, `"$1": "$2",`)
				bodyStr = strings.Replace(bodyStr, `90"LAR`, `90 LAR`, -1)

				//Modification de la réponse
				//log.Printf("Username = %s\n", Username)
				//log.Printf("Password = %s\n", Password)

				/*
					config.Xtream.Hostname = gjson.Get(bodyStr, "server_info.url").String()
					port, _ := strconv.Atoi(gjson.Get(bodyStr, "server_info.port").String())
					config.Xtream.Port = uint16(port)
				*/

				isT := "0"
				if usr.IsTrial {
					isT = "1"
				}
				bodyStr, _ = sjson.Set(bodyStr, "user_info.username", usr.Username)
				bodyStr, _ = sjson.Set(bodyStr, "user_info.password", usr.Password)

				bodyStr, _ = sjson.Set(bodyStr, "user_info.status", usr.Status)
				bodyStr, _ = sjson.Set(bodyStr, "user_info.exp_date", strconv.Itoa(int(usr.ExpDate.Unix())))
				bodyStr, _ = sjson.Set(bodyStr, "user_info.is_trial", isT)
				bodyStr, _ = sjson.Set(bodyStr, "user_info.created_at", strconv.Itoa(int(usr.CreatedAt.Unix())))
				bodyStr, _ = sjson.Set(bodyStr, "user_info.max_connections", strconv.Itoa(usr.MaxConnections))

				bodyStr, _ = sjson.Set(bodyStr, "server_info.url", config.Server.Hostname)
				bodyStr, _ = sjson.Set(bodyStr, "server_info.port", strconv.Itoa(int(config.Server.Port)))
				bodyStr, _ = sjson.Delete(bodyStr, "server_info.https_port")
				bodyStr, _ = sjson.Set(bodyStr, "server_info.server_protocol", "http")

				//Fin modification de la réponse
			}

			fmt.Println(urlRequest)
		}
		//db.GetUser(Username, Password)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(bodyStr))
	}
}
