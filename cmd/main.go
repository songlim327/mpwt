package main

import (
	"fmt"
	"mpwt/internal/config"
	"mpwt/internal/core"
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
		file, err := os.OpenFile("./mpwt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Failed to create log file: %v", err))
		}
		defer file.Close()

		log.NewLogWithFile(log.EnvProduction, file)
	}

	log.Debug("halo")
	log.Info("halo")

	terminalConf := &core.TerminalConfig{
		Maximize:        conf.Maximize,
		Direction:       conf.Direction,
		Columns:         conf.Columns,
		OpenInNewWindow: conf.OpenInNewWindow,
		Commands:        []string{"echo 1", "echo 2", "echo 3", "echo 4", "echo 5", "echo 6", "echo 7", "echo 8", "echo 9", "echo 10"},
	}
	err = core.OpenWt(terminalConf)
	if err != nil {
		log.Error(err)
	}

	// core.InitTea()

	// sizes, err := core.CalculatePaneSize(5)
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Debug(sizes)

	// cmd := exec.Command("wt")
	// cmd.Run()
}
