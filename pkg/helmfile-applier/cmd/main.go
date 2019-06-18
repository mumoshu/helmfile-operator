package main

import (
	"fmt"
	"github.com/mumoshu/helmfile-operator/pkg/helmfile-applier"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	name := "helmfile-applier"

	cmd := &cobra.Command{
		Use: fmt.Sprintf("%s [flags]", name),
	}

	flagset := cmd.PersistentFlags()

	var source string
	var once bool
	var environment string

	flagset.StringVarP(&source, "file", "f", "", "desired state file to be used to reconcile the cluster. currently only helmfile-style state file is supported.")
	flagset.StringVarP(&environment, "environment", "e", "", "helmfile environment name to be specified for deployments")
	flagset.BoolVar(&once, "once", false, "run once and exit immediately. primarily for testing and development purpose")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		r, err := helmfile_applier.New(
			nil,
			helmfile_applier.Source(source),
			helmfile_applier.Once(once),
			helmfile_applier.Environment(environment),
		)
		if err != nil {
			return err
		}

		return r.Run()
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s exiting: %v", name, err)
		os.Exit(1)
	}
}
