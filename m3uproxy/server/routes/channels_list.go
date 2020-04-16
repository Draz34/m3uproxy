package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/Draz34/m3uproxy/server/webutils"
)

const UriChannelList = "/channels"

type LoadingChannelsError struct {
	Msg    string
	Error  error
	Status int
}

func ChannelListRouter(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request)) {
	return UriChannelList, func(w http.ResponseWriter, r *http.Request) {
		bytes, err := LoadList(config)

		if err != nil {
			switch err.Status {
			case http.StatusBadGateway:
				webutils.BadGateway(err.Msg, err.Error, w)
				return
			default:
				webutils.InternalServerError(err.Msg, err.Error, w)
				return
			}
		}

		webutils.Success(bytes, w)
	}
}

func LoadList(config *config.Config) ([]byte, *LoadingChannelsError) {
	resp, err := http.Get(config.M3u.Url)
	if err != nil {
		return nil, &LoadingChannelsError{
			Msg:    "Error loading channels list",
			Error:  err,
			Status: http.StatusBadGateway,
		}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &LoadingChannelsError{
			Msg:    "Error loading body from response",
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, &LoadingChannelsError{
			Msg:    "Error closing body response",
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	if resp.StatusCode == 200 {
		db.Reset()

		b, err = modifyResponse(config, string(b))
		if err != nil {
			return nil, &LoadingChannelsError{
				Msg:    "Error modifying response",
				Error:  err,
				Status: http.StatusInternalServerError,
			}
		}

	}

	return b, nil
}

func modifyResponse(config *config.Config, payload string) ([]byte, error) {
	newPayload1 := strings.Replace(payload, "\r\n", "\n", -1)
	newPayload2 := strings.Replace(newPayload1, "\n", "", -1)
	newPayload3 := strings.Replace(newPayload2, "http://", "\nhttp://", -1)
	newPayload4 := strings.Replace(newPayload3, "#EXTINF", "\n#EXTINF", -1)
	lines := strings.Split(string(newPayload4), "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Line holding metadata. Something like
		// #EXTINF:-1 tvg-id="ABC HD" tvg-name="ABC FHD" ...
		if !strings.HasPrefix(line, "http://") {
			continue
		}

		// line is a channel address.
		// Override channel address with proxyHost address
		channel, err := db.RegisterChannel(line)
		if err != nil {
			return nil, fmt.Errorf("error registering m3u url. %v", err)
		}

		lines[i] = GetChannelUrl(config, channel.Id)
	}

	return []byte(strings.Join(lines, "\n")), nil
}
