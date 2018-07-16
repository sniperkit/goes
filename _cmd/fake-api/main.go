package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/sniperkit/janus/config"
	"github.com/sniperkit/janus/server"
)

const VERSION = "1.3.0"

var (
	currentWorkDir, _ = os.Getwd()
	configPrefixPath  = flag.String("config-dir", currentWorkDir, "Config prefix path")
	configFilename    = flag.String("config-file", "config.json", "Config filename")
	resDefaultDir     = filepath.Join(currentWorkDir, "data")
	resPrefixPath     = flag.String("resource-dir", resDefaultDir, "Resources prefix path")
)

func main() {
	fmt.Printf("Janus - fake rest api server (%s) \n", VERSION)
	flag.Parse()

	c := getConfig(*configPrefixPath, *configFilename, *resPrefixPath)

	go server.Start(c)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	server.Stop()
}

// getConfig get the configuration from the config file.
func getConfig(configPrefixPath, configFilename, resPrefixPath string) *config.Config {
	configFile := filepath.Join(configPrefixPath, configFilename)
	c, err := config.ParseFile(configFile)
	if err != nil {
		fmt.Printf("Config prefix path: %s\n", configPrefixPath)
		fmt.Printf("Config file name: %s\n", configFilename)
		fmt.Printf("Config file path: %s\n", configFile)
		fmt.Printf("Loading config error: %s\n", err.Error())
		os.Exit(1)
	}

	c.Path = resPrefixPath
	return c
}
