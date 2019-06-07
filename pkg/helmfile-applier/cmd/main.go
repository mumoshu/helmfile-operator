package main

import (
	"fmt"
	"github.com/mumoshu/appliance-operator/pkg/helmfile-applier"
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
	var environment string

	flagset.StringVarP(&source, "file", "f", "", "desired state file to be used to reconcile the cluster. currently only helmfile-style state file is supported.")
	flagset.StringVarP(&environment, "environment", "e", "", "helmfile environment name to be specified for deployments")
	flagset.BoolVar(&once, "once", false, "run once and exit immediately. primarily for testing and development purpose")

	r, err := helmfile_applier.New(
		nil,
		helmfile_applier.Source(source),
		helmfile_applier.Once(once),
		helmfile_applier.Environment(environment),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	if err := r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "genericcontroller exiting: %v", err)
		os.Exit(1)
	}
}
