package main

import (
	"fmt"
	"mpwt/internal/config"
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

	// sizes, err := core.CalculatePaneSize(5)
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Debug(sizes)

	// cmd := exec.Command("wt")
	// cmd.Run()
}
