package webutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Draz34/m3uproxy/db"
)

func Success(b []byte, w http.ResponseWriter) {
	w.WriteHeader(200)
	writePayload(b, w, false)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
}

func BadRequest(msg string, cause error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusBadRequest)
	writePayload([]byte(msg+"\n"+cause.Error()), w, true)
}

func InternalServerError(msg string, cause error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	writePayload([]byte(msg+"\n"+cause.Error()), w, true)
}

func BadGateway(msg string, cause error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusBadGateway)
	writePayload([]byte(msg+"\n"+cause.Error()), w, true)
}

func writePayload(payload []byte, w http.ResponseWriter, isError bool) {
	if isError {
		log.Printf("An error occured: %s\n", payload)
	}

	_, err := w.Write(payload)
	if err != nil {
		log.Printf("Error writing content to http.ResponseWriter: payload=%s, err=%v", payload, err)
	}
}

func TracingRedirect(myURL string) {
	nextURL := myURL
	var i int
	for i < 100 {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}

		resp, err := client.Get(nextURL)

		if err != nil {
			//fmt.Println(err)
		}

		uHost, _ := url.Parse(nextURL)

		hostA := strings.Split(uHost.Host, ":")
		host := ""
		port := ""
		if len(hostA) > 1 {
			host = hostA[0]
			port = hostA[1]
		}
		paths := strings.Split(uHost.Path, "/")
		username := ""
		password := ""
		if len(paths) > 2 {
			username = paths[2]
			password = paths[3]
		}

		var p = db.XtreamProxy{
			Domain:   host,
			Port:     port,
			Username: username,
			Password: password,
			Md5:      GetMD5Hash(host + port + username + password),
			Url:      uHost.String(),
		}

		db.CreateXtreamProxy(p)

		if resp.StatusCode == 200 {
			fmt.Println("Done!")
			break
		} else {
			fmt.Println("StatusCode:", resp.StatusCode)
			fmt.Println(resp.Request.URL)

			nextURL = resp.Header.Get("Location")
			i += 1
		}
	}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
