package main

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/mumoshu/helmfile-server/pkg/genericcontroller"
	"os"
)

func main() {
	assetsDir := "assets"

	// The second argument to packr.New must be a local variable or a string literal
	// in order for `packr2 build` to successfully determine the directory to be packed
	box := packr.New("Bundled Addon Assets", assetsDir)

	r, err := genericcontroller.New(
		box,
		genericcontroller.AssetDir(assetsDir),
		genericcontroller.Source("assets/helmfile.yaml"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	if err := r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
