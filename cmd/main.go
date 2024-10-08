//go:generate goversioninfo -platform-specific=true assets/versioninfo.json
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
	"path/filepath"
)

func main() {
	// Identify application enviroment (development/production)
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Get executable path
	exeDir, err := getExecDirectory()
	if err != nil {
		panic(fmt.Errorf("failed to get executable directory: %v", err))
	}

	// Intialize logger and set config file path based on application environment
	configPath := ""
	if *debug {
		configPath = filepath.Join(exeDir, "/config/config.dev.yaml")
		log.NewLog(log.EnvDevelopment)
	} else {
		configPath = filepath.Join(exeDir, "/config.yaml")
		file, err := os.OpenFile(exeDir+"/mpwt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("failed to create log file: %v", err))
		}
		defer file.Close()

		log.NewLogWithFile(log.EnvProduction, file)
	}

	// Read config from yaml config file
	mgr := config.NewConfigManager(configPath)
	conf, err := mgr.NewConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read config file: %v", err))
	}

	// Initialize database connection
	r, err := repository.NewDbConn(exeDir + "/mpwt.db")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize sqlite: %v", err))
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
		ConfigMgr:  mgr,
	}

	// Start terminal application
	err = tui.InitTea(tuiConf)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to run tui: %v", err))
	}
}

// getExecDirectory returns the directory containing the application executable
func getExecDirectory() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve executable path: %v", err)
	}

	// Evaluate symlinks to prevent unstable result from os.Executable
	resolvedPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlink: %v", err)
	}

	// Get the directory of the resolved executable
	return filepath.Dir(resolvedPath), nil
}
