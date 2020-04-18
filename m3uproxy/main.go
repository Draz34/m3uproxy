package main

import (
	"flag"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/db"
	"github.com/Draz34/m3uproxy/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var m3uServerConfig *config.Config

	var ymlConfigFile string
	flag.StringVar(&ymlConfigFile, "file", "", "Configuration file")
	flag.Parse()

	if ymlConfigFile != "" {
		m3uServerConfig = config.LoadYml(ymlConfigFile)

	} else {
		m3uServerConfig = config.LoadEnv()
	}
	db.Test()
	config.Validate(m3uServerConfig)
	server.Start(m3uServerConfig)
}
