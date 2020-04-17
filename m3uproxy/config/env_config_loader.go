package config

import (
	"log"
	"os"
	"strconv"
)

const (
	M3uProxyPort          = "M3U_PROXY_PORT"
	M3uProxyHostname      = "M3U_PROXY_HOSTNAME"
	M3uProxyXtremPort     = "M3U_PROXY_XTREAM_PORT"
	M3uProxyXtremHostname = "M3U_PROXY_XTREAM_HOSTNAME"
	M3uProxyXtremSurname  = "M3U_PROXY_XTREAM_SURNAME"
	M3uProxyXtremUsername = "M3U_PROXY_XTREAM_USERNAME"
	M3uProxyXtremPassword = "M3U_PROXY_XTREAM_PASSWORD"
	M3uProxyM3uUrl        = "M3U_PROXY_CHANNELS_URL"
)

func LoadEnv() *Config {
	var config = &Config{}

	config.Server.Port = 9090
	config.Server.Hostname = "localhost"

	config.Xtream.Port = 7713
	config.Xtream.Hostname = "10.10.10.10"
	config.Xtream.Surname = "toto"
	config.Xtream.Username = "root"
	config.Xtream.Password = "password"

	port := os.Getenv(M3uProxyPort)
	if port != "" {
		envPort, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			log.Fatalf("Error parsing server port number: %s", port)
		}

		config.Server.Port = uint16(envPort)
	}

	envHostname := os.Getenv(M3uProxyHostname)
	if envHostname != "" {
		config.Server.Hostname = envHostname
	}

	XtreamPort := os.Getenv(M3uProxyXtremPort)
	if XtreamPort != "" {
		envXtreamPort, err := strconv.ParseInt(XtreamPort, 10, 64)
		if err != nil {
			log.Fatalf("Error parsing server port number: %s", XtreamPort)
		}

		config.Xtream.Port = uint16(envXtreamPort)
	}

	envXtreamHostname := os.Getenv(M3uProxyXtremHostname)
	if envXtreamHostname != "" {
		config.Xtream.Hostname = envXtreamHostname
	}

	envXtreamSurname := os.Getenv(M3uProxyXtremSurname)
	if envXtreamSurname != "" {
		config.Xtream.Surname = envXtreamSurname
	}

	envXtreamUsername := os.Getenv(M3uProxyXtremUsername)
	if envXtreamUsername != "" {
		config.Xtream.Username = envXtreamHostname
	}

	envXtreamPassword := os.Getenv(M3uProxyXtremPassword)
	if envXtreamPassword != "" {
		config.Xtream.Password = envXtreamPassword
	}

	config.M3u.Url = os.Getenv(M3uProxyM3uUrl)
	return config
}
