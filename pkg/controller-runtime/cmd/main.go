package main

import (
	"flag"
	"fmt"
	"github.com/mumoshu/helmfile-operator/pkg/controller-runtime"
	"os"
)

func main() {
	var name = flag.String("name", "helmfile-operator", "Operator name included in log messages")
	var configPath = flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if err := controller_runtime.Run(*name, *configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
