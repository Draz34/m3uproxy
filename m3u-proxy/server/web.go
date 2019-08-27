package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hmarcelino/m3u-proxy/config"
	"github.com/hmarcelino/m3u-proxy/server/routes"
	"log"
	"net/http"
)

const Logo = ` 
___  ___ _____ _   _     ____________ _______   ____   __
|  \/  ||____ | | | |    | ___ \ ___ \  _  \ \ / /\ \ / /
| .  . |    / / | | |    | |_/ / |_/ / | | |\ V /  \ V / 
| |\/| |    \ \ | | |    |  __/|    /| | | |/   \   \ /  
| |  | |.___/ / |_| |    | |   | |\ \\ \_/ / /^\ \  | |  
\_|  |_/\____/ \___/     \_|   \_| \_|\___/\/   \/  \_/

is accepting requests in port :%d
* http://127.0.0.1:%d'
* http://%s:%d'

`

func Start(config *config.Config) {
	muxRouter := mux.NewRouter()

	register(muxRouter, config, routes.RootRouter)
	register(muxRouter, config, routes.PingRouter)
	register(muxRouter, config, routes.ChannelListRouter)
	register(muxRouter, config, routes.ChannelRoute)

	fmt.Printf(
		Logo,
		config.Server.Port,
		config.Server.Port,
		config.Server.Hostname,
		config.Server.Port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), muxRouter)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func register(mux *mux.Router, config *config.Config, route func(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request))) {
	path, handler := route(config)
	mux.HandleFunc(path, handler)
}
