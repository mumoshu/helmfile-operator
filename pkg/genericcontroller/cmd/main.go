package main

import (
	"fmt"
	"github.com/mumoshu/helmfile-server/pkg/genericcontroller"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use: "genericcontroller [flags]",
	}

	flagset := cmd.PersistentFlags()

	var source string
	var once bool

	flagset.StringVarP(&source, "file", "f", "", "desired state file to be used to reconcile the cluster. currently only helmfile-style state file is supported.")
	flagset.BoolVar(&once, "once", false, "run once and exit immediately. primarily for testing and development purpose")

	r, err := genericcontroller.New(nil, genericcontroller.Source(source), genericcontroller.Once(once))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	if err := r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "genericcontroller existing: %v", err)
		os.Exit(1)
	}
}
