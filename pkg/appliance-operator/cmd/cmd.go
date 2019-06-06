package cmd

import (
	"flag"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/mumoshu/helmfile-server/pkg/apputil"
	"github.com/mumoshu/helmfile-server/pkg/helmfile-operator"
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

	out, err := apputil.RunCommand("kubectl", "apply", "-f", "assets/init")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Apply sucecssful: %s\n", out)

	if err := helmfile_operator.Run(*name, *configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
