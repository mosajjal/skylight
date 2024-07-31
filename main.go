package main

import (
	_ "embed"
	"flag"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	skylight "github.com/mosajjal/skylight/pkg"
	"github.com/phuslu/log"
)

var nocolorLog = strings.ToLower(os.Getenv("NO_COLOR")) == "true"

//go:embed config.hcl
var defaultConfig []byte

var (
	cfgFile          = flag.String("config", "", "path to the config file")
	printDefaultConf = flag.Bool("defaultconfig", false, "write the default config to stdout")
)

func main() {

	flag.Parse()

	var config skylight.Config
	if err := hclsimple.DecodeFile(*cfgFile, nil, &config); err != nil {
		log.Fatal().Msgf("Failed to load configuration: %s", err)
	}

	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Caller:     1,
			Level:      log.ParseLevel(config.LogLevel),
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

	// load the default config
	if *printDefaultConf {
		os.Stdout.Write(defaultConfig)
		os.Exit(0)
	}

	// this run is non-blocking
	skylight.Run(config)

	// wait forever
	// TODO: add a signal handler to gracefully shutdown
	select {}
}
