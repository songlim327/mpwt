package main

import (
	"fmt"
	"mpwt/internal/config"
	"mpwt/internal/core"
	"mpwt/internal/repository"
	"mpwt/internal/tui"
	"mpwt/pkg/log"
	"os"
)

func main() {
	// Read config from yaml config file
	pwd, _ := os.Getwd()
	conf, err := config.NewConfig(pwd + "/internal/config/config.dev.yaml")
	if err != nil {
		panic(err)
	}

	// Intialize logger
	if conf.Debug {
		log.NewLog(log.EnvDevelopment)
	} else {
		file, err := os.OpenFile(pwd+"/mpwt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Failed to create log file: %v", err))
		}
		defer file.Close()

		log.NewLogWithFile(log.EnvProduction, file)
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
