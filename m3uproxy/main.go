package main

import (
	"flag"

	"github.com/Draz34/m3uproxy/config"
	"github.com/Draz34/m3uproxy/server"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	var m3uServerConfig *config.Config

	var ymlConfigFile string
	flag.StringVar(&ymlConfigFile, "file", "", "Configuration file")
	flag.Parse()

	confType := 1
	if ymlConfigFile != "" {
		m3uServerConfig = config.LoadYml(ymlConfigFile)

	} else {
		m3uServerConfig = config.LoadEnv()
		confType = 2
	}

	config.Validate(m3uServerConfig)
	server.Start(m3uServerConfig, confType)
}
