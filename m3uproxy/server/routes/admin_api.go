package routes

import (
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

		urlRequest := "/admin_api.php?action=" + action + "&username=" + u.Username + "&password=" + u.Password + "&status=" + u.Status + "&exp_date=" + strconv.Itoa(int(u.ExpDate.Unix())) + "&is_trial=" + r.FormValue("is_trial") + "&created_at=" + strconv.Itoa(int(u.CreatedAt.Unix())) + "&max_connections=" + string(u.MaxConnections)
		fmt.Println(urlRequest)

		w.Header().Set("Content-Type", "application/json")
		webutils.Success([]byte(`[{"id":"1","username":"username1","credits":"0","group_id":"1","group_name":"Administrators","last_login":"1511883014","date_registered":"1421604973","email":"test1@my-email.com","ip":"8.8.8.8","status":"1"},{"id":"31","username":"username2","credits":"0","group_id":"1","group_name":"Administrators","last_login":"1539813642","date_registered":"1473632400","email":"test2@my-email.com","ip":"8.8.8.8","status":"1"}]`), w)
	}
}
