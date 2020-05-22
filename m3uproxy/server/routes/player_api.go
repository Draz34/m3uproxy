package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
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

		if categorieNum == "ALL" {
			categorieNum = ""
		}

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

			bodyStr = string(body)

			//fmt.Println(string(body))

			urlRequest := urlString
			if Action != "" {
				var sortDatas bool = true

				urlRequest = urlString + "?action=" + Action + "&category_id=" + categorieNum + "&stream_id=" + streamNum + "&series_id=" + serieNum + "&vod_id=" + vodNum + "&limit=" + Limit
				switch Action {
				case "get_vod_categories":
					var moviesCategories = make([]db.MovieCategorie, 0)
					set := make(map[string]db.MovieCategorie)

					err = json.Unmarshal(body, &moviesCategories)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}

					for k, _ := range moviesCategories {
						set[moviesCategories[k].CategoryID] = moviesCategories[k]
						moviesCategories[k].CategoryName = strings.TrimSpace(moviesCategories[k].CategoryName)
					}

					sort.SliceStable(moviesCategories, func(i, j int) bool {
						return moviesCategories[i].CategoryName < moviesCategories[j].CategoryName
					})

					//fmt.Println(moviesCategories2)

					//Liste les vod, rajoute les catégories si elles n'existent pas
					var movies = make([]db.Movie, 0)
					formData = url.Values{
						"username": {config.Xtream.Username},
						"password": {config.Xtream.Password},
						"action":   {"get_vod_streams"},
					}

					resp2, err := http.PostForm(urlString, formData)
					if err != nil {
						print(err)
					}

					defer resp2.Body.Close()
					body2, err := ioutil.ReadAll(resp2.Body)
					if err != nil {
						print(err)
					}

					err = json.Unmarshal(body2, &movies)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}

					for k, _ := range movies {
						if _, ok := set[movies[k].CategoryID]; ok {
							//fmt.Println("element found")
						} else {
							newCat := db.MovieCategorie{
								CategoryID:   movies[k].CategoryID,
								CategoryName: "Categorie " + movies[k].CategoryID,
								ParentID:     0,
							}
							set[movies[k].CategoryID] = newCat
							moviesCategories = append(moviesCategories, newCat)
						}
					}
					//Fin liste vod

					body, err = json.Marshal(moviesCategories)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}
					if sortDatas {
						bodyStr = string(body)
					}
				case "get_vod_streams":
					var movies = make([]db.Movie, 0)
					err = json.Unmarshal(body, &movies)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}

					for k, _ := range movies {
						movies[k].Name = strings.TrimSpace(movies[k].Name)
					}

					sort.SliceStable(movies, func(i, j int) bool {
						return movies[i].Name < movies[j].Name
					})

					//fmt.Println(movies)

					body, err = json.Marshal(movies)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}
					if sortDatas {
						bodyStr = string(body)
					}
				case "get_series":
					var series = make([]db.Serie, 0)
					err = json.Unmarshal(body, &series)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}

					for k, _ := range series {
						series[k].Name = strings.TrimSpace(series[k].Name)
					}

					sort.SliceStable(series, func(i, j int) bool {
						return series[i].Name < series[j].Name
					})

					//fmt.Println(movies)

					body, err = json.Marshal(series)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}
					if sortDatas {
						bodyStr = string(body)
					}
				case "get_live_streams":
					var lives = make([]db.Live, 0)
					err = json.Unmarshal(body, &lives)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}

					for k, _ := range lives {
						lives[k].Name = strings.TrimSpace(lives[k].Name)
					}

					sort.SliceStable(lives, func(i, j int) bool {
						return lives[i].Name < lives[j].Name
					})

					//fmt.Println(movies)

					body, err = json.Marshal(lives)
					if err != nil {
						fmt.Print(err)
						sortDatas = false
					}
					if sortDatas {
						bodyStr = string(body)
					}
				}
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

func addMovie(proxy db.XtreamProxy) (version string, liveCount int, movieCount int, serieCount int) {
	urlString := "http://" + proxy.Domain + ":" + proxy.Port + "/player_api.php"
	urlString2 := "http://" + proxy.Domain + ":" + proxy.Port + "/panel_api.php"

	//Lives
	formData := url.Values{
		"username": {proxy.Username},
		"password": {proxy.Password},
		"action":   {"get_live_streams"},
	}

	version = "2"
	resp2, err := http.PostForm(urlString, formData)
	if err != nil {
		print(err)
		urlString = urlString2
		resp2, err = http.PostForm(urlString, formData)
		if err != nil {
			print(err)
		}
		version = "1"
	}

	defer resp2.Body.Close()
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		print(err)
	}

	var jsonObjs interface{}
	json.Unmarshal(body2, &jsonObjs)
	objSlice, ok := jsonObjs.([]interface{})

	if !ok {
		fmt.Println("cannot convert the JSON objects")
	}
	liveCount = len(objSlice)

	//Movies
	formData = url.Values{
		"username": {proxy.Username},
		"password": {proxy.Password},
		"action":   {"get_vod_streams"},
	}

	resp2, err = http.PostForm(urlString, formData)
	if err != nil {
		print(err)
	}

	defer resp2.Body.Close()
	body2, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		print(err)
	}

	json.Unmarshal(body2, &jsonObjs)
	objSlice, ok = jsonObjs.([]interface{})

	if !ok {
		fmt.Println("cannot convert the JSON objects")
	}
	movieCount = len(objSlice)

	//Series
	formData = url.Values{
		"username": {proxy.Username},
		"password": {proxy.Password},
		"action":   {"get_series"},
	}

	resp2, err = http.PostForm(urlString, formData)
	if err != nil {
		print(err)
	}

	defer resp2.Body.Close()
	body2, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		print(err)
	}

	json.Unmarshal(body2, &jsonObjs)
	objSlice, ok = jsonObjs.([]interface{})

	if !ok {
		fmt.Println("cannot convert the JSON objects")
	}
	serieCount = len(objSlice)

	return
}
