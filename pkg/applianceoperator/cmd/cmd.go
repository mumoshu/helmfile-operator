package cmd

import (
	"flag"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/mumoshu/helmfile-server/pkg/apputil"
	"github.com/mumoshu/helmfile-server/pkg/genericoperator"
	"os"
)

func Run() {
	assetsDir := "assets"

	// The second argument to packr.New must be a local variable or a string literal
	// in order for `packr2 build` to successfully determine the directory to be packed
	box := packr.New("Bundled Appliance Assets", assetsDir)

	l := apputil.NewLogger(os.Stderr, "debug")
	syncer, err := apputil.New(
		apputil.Box(box),
		apputil.Logger(l),
		apputil.Assets(assetsDir),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var name = flag.String("name", "appliance-operator", "Operator name included in log messages")
	var configPath = flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if err := syncer.SyncOnce(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := genericoperator.Run(*name, *configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
