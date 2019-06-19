package main

import (
	"encoding/json"
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
	var valuesFiles []string
	var valuesJson string
	var helmx bool

	flagset.StringVarP(&source, "file", "f", "", "desired state file to be used to reconcile the cluster. currently only helmfile-style state file is supported.")
	flagset.StringVarP(&environment, "environment", "e", "", "helmfile environment name to be specified for deployments")
	flagset.StringSliceVar(&valuesFiles, "valuesFile", []string{}, "values files to be passed to helmfile")
	flagset.StringVar(&valuesJson, "values", "{}", "values to be passed to helmfile, as JSON object")
	flagset.BoolVar(&helmx, "helm-x", false, "enable helm-x integration. Required for kustomize, k8s manifests, patching support")
	flagset.BoolVar(&once, "once", false, "run once and exit immediately. primarily for testing and development purpose")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(valuesJson), &m); err != nil {
			return err
		}

		r, err := helmfile_applier.New(
			nil,
			helmfile_applier.Source(source),
			helmfile_applier.Once(once),
			helmfile_applier.Environment(environment),
			helmfile_applier.Values(m),
			helmfile_applier.ValuesFiles(valuesFiles),
			helmfile_applier.HelmX(helmx),
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
