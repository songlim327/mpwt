package main

import (
	"flag"
	"fmt"
	"mpwt/internal/config"
	"mpwt/internal/core"
	"mpwt/internal/repository"
	"mpwt/internal/tui"
	"mpwt/pkg/log"
	"os"
)

func main() {
	// Identify application enviroment (development/production)
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Intialize logger and set config file path based on application environment
	pwd, _ := os.Getwd()
	configPath := pwd
	if *debug {
		configPath += "/config/config.dev.yaml"
		log.NewLog(log.EnvDevelopment)
	} else {
		configPath += "/config.yaml"
		file, err := os.OpenFile(pwd+"/mpwt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Failed to create log file: %v", err))
		}
		defer file.Close()

		log.NewLogWithFile(log.EnvProduction, file)
	}

	// Read config from yaml config file
	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	// Initialize database connection
	r, err := repository.NewDbConn(pwd + "/mpwt.db")
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	// Initialize tui configuration
	tuiConf := &tui.TuiConfig{
		TerminalConfig: &core.TerminalConfig{
			Maximize:     conf.Maximize,
			Direction:    conf.Direction,
			Columns:      conf.Columns,
			OpenInNewTab: conf.OpenInNewTab,
		},
		Repository: r,
	}

	// Start terminal application
	err = tui.InitTea(tuiConf)
	if err != nil {
		log.Fatal(err)
	}
}
