package main

import (
	"flag"
	"fmt"
	"github.com/mumoshu/helmfile-operator/pkg/controller-runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
)

func main() {
	var name = flag.String("name", "helmfile-operator", "Operator name included in log messages")
	var configPath = flag.String("config", "", "Path to configuration file")

	flag.Parse()

	resource := schema.GroupVersionKind{
		Group:   "apps.mumoshu.github.io",
		Kind:    "Appliance",
		Version: "v1alpha1",
	}
	if err := controller_runtime.Run(*name, resource, controller_runtime.Conf(*configPath)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
