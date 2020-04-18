package routes

import (
	"fmt"
	"net/http"

	"github.com/Draz34/m3uproxy/config"
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
		username := r.FormValue("username")
		password := r.FormValue("password")
		status := r.FormValue("status")
		expDate := r.FormValue("exp_date")
		isTrial := r.FormValue("is_trial")
		createdAt := r.FormValue("created_at")
		maxConnections := r.FormValue("max_connections")

		urlRequest := "/admin_api.php?action=" + action + "&username=" + username + "&password=" + password + "&status=" + status + "&exp_date=" + expDate + "&is_trial=" + isTrial + "&created_at=" + createdAt + "&max_connections=" + maxConnections
		fmt.Println(urlRequest)

		w.Header().Set("Content-Type", "application/json")
		webutils.Success([]byte(`[{"id":"1","username":"username1","credits":"0","group_id":"1","group_name":"Administrators","last_login":"1511883014","date_registered":"1421604973","email":"test1@my-email.com","ip":"8.8.8.8","status":"1"},{"id":"31","username":"username2","credits":"0","group_id":"1","group_name":"Administrators","last_login":"1539813642","date_registered":"1473632400","email":"test2@my-email.com","ip":"8.8.8.8","status":"1"}]`), w)
	}
}
