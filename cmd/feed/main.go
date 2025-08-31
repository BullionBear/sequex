package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BullionBear/sequex/internal/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	// Setup logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return "", filename + ":" + string(rune(f.Line))
		},
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Parse config file
	var configFile string
	flag.StringVar(&configFile, "c", "./config.yml", "Path to configuration file")
	flag.Parse()
	log.Infof("Starting feed service with config file: %s", configFile)
	// Read config file
	config := &config.FeedConfig{}
	if data, err := os.ReadFile(configFile); err == nil {
		if err := yaml.Unmarshal(data, config); err != nil {
			log.Fatalf("Failed to parse config file: %v", err)
		}
	} else {
		log.Warnf("Could not read config file %s: %v", configFile)
	}

	log.Infof("Starting feed service with config: %+v", config)

}
