package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/Draz34/m3uproxy/server/webutils"
)

type BasicAuthFunc func(username, password string) bool

func (f BasicAuthFunc) RequireAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
	w.WriteHeader(401)
}

func (f BasicAuthFunc) Authenticate(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	return ok && f(username, password)
}

func SimpleBasicAuth(username, password string) BasicAuthFunc {
	return BasicAuthFunc(func(user, pass string) bool {
		return username == user && password == pass
	})
}

func AdminApiRoute(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {
	return "/admin_api.php", func(w http.ResponseWriter, r *http.Request) {

		f := SimpleBasicAuth(config.Server.AdminLogin, config.Server.AdminPassword)

		if !f.Authenticate(r) {
			f.RequireAuth(w)
			return
		}

		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}

		action := r.FormValue("action")

		var jsonResponse []byte

		switch action {
		case "create_user":
			dateExp, err := time.Parse("2006-01-02 15:04:05", r.FormValue("exp_date"))
			if err != nil {
				fmt.Println(err.Error())
			}

			dateCrea, err := time.Parse("2006-01-02 15:04:05", r.FormValue("created_at"))
			if err != nil {
				fmt.Println(err.Error())
			}

			isT, err := strconv.ParseBool(r.FormValue("is_trial"))
			if err != nil {
				fmt.Println(err.Error())
			}

			maxCo, err := strconv.ParseInt(r.FormValue("max_connections"), 10, 32)
			if err != nil {
				fmt.Println(err.Error())
			}

			var u = db.User{
				Username:       r.FormValue("username"),
				Password:       r.FormValue("password"),
				Status:         r.FormValue("status"),
				ExpDate:        dateExp,
				IsTrial:        isT,
				CreatedAt:      dateCrea,
				MaxConnections: int(maxCo),
			}

			db.CreateUser(u)

			jsonResponse, err = json.Marshal(u)
			if err != nil {
				webutils.InternalServerError("Error building jsonResponse from a User", err, w)
			}

			urlRequest := "/admin_api.php?action=" + action + "&username=" + u.Username + "&password=" + u.Password + "&status=" + u.Status + "&exp_date=" + strconv.Itoa(int(u.ExpDate.Unix())) + "&is_trial=" + r.FormValue("is_trial") + "&created_at=" + strconv.Itoa(int(u.CreatedAt.Unix())) + "&max_connections=" + string(u.MaxConnections)
			fmt.Println(urlRequest)
		case "users_list":
			var err error
			jsonResponse, err = json.Marshal(db.GetAllUser())
			if err != nil {
				webutils.InternalServerError("Error building jsonResponse from a Xtream Proxy", err, w)
			}
		case "servers_list":
			var err error
			jsonResponse, err = json.Marshal(db.GetAllXtreamProxy())
			if err != nil {
				webutils.InternalServerError("Error building jsonResponse from a Xtream Proxy", err, w)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		webutils.Success(jsonResponse, w)
	}
}
