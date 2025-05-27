package main

import (
	"flag"
	"fmt"
)

func main() {
	configFilename := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config := MustLoadConfig(*configFilename)
	ctx := AppCtx{
		Config: config,
	}
	fmt.Print(ctx)
}
