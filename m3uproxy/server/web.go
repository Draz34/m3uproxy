package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/Draz34/m3uproxy/server/routes"
	"github.com/gorilla/mux"
)

const Logo = ` 
___  ___ _____ _   _     ____________ _______   ____   __
|  \/  ||____ | | | |    | ___ \ ___ \  _  \ \ / /\ \ / /
| .  . |    / / | | |    | |_/ / |_/ / | | |\ V /  \ V / 
| |\/| |    \ \ | | |    |  __/|    /| | | |/   \   \ /  
| |  | |.___/ / |_| |    | |   | |\ \\ \_/ / /^\ \  | |  
\_|  |_/\____/ \___/     \_|   \_| \_|\___/\/   \/  \_/

is accepting requests in port :%d
* http://127.0.0.1:%d
* http://%s:%d

xtream config :
* url : http://%s:%d
* surname : %s
* username : %s
* password : %s

`

func Start(config *config.Config) {
	muxRouter := mux.NewRouter()

	register(muxRouter, config, routes.RootRouter)
	register(muxRouter, config, routes.PingRouter)
	register(muxRouter, config, routes.ChannelListRouter)
	register(muxRouter, config, routes.ChannelRoute)
	register(muxRouter, config, routes.ChannelInfoRoute)
	register(muxRouter, config, routes.PanelApiRoute)
	register(muxRouter, config, routes.LiveRoute)

	//Log not found routes
	//muxRouter.NotFoundHandler = muxRouter.NewRoute().HandlerFunc(http.NotFound).GetHandler()
	//muxRouter.Use(simpleMw)
	muxRouter.NotFoundHandler = muxRouter.NewRoute().HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("404 not found"))
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		}).GetHandler()

	fmt.Printf(
		Logo,
		config.Server.Port,
		config.Server.Port,
		config.Server.Hostname,
		config.Server.Port,
		config.Xtream.Hostname,
		config.Xtream.Port,
		config.Xtream.Surname,
		config.Xtream.Username,
		config.Xtream.Password)

	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Server.Port), Handler: muxRouter}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			//log.Fatalf("Error starting server: %v", err)
		}
	}()

	_, err := routes.LoadList(config)
	if routes.LoadList(config); err != nil {
		log.Fatalf(err.Msg+" %v", err.Error)
	}

	log.Println("List loaded successfully with " + strconv.Itoa(db.ChannelsLen()) + " url(s)")

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		// ignoring error
	}
}

func register(mux *mux.Router, config *config.Config, route func(config *config.Config) (string, func(w http.ResponseWriter, r *http.Request))) {
	path, handler := route(config)
	mux.HandleFunc(path, handler)
}

func simpleMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
