package cmd

import (
	"flag"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/mumoshu/helmfile-operator/pkg/apputil"
	"github.com/mumoshu/helmfile-operator/pkg/controller-runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
)

func Run() {
	assetsDir := "assets"

	// The second argument to packr.New must be a local variable or a string literal
	// in order for `packr2 build` to successfully determine the directory to be packed
	box := packr.New("Bundled Helmfile Assets", assetsDir)

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

	var name = flag.String("name", "helmfile-operator", "Operator name included in log messages")
	var configPath = flag.String("config", "", "Path to configuration file")
	applierImage := flag.String("image", "", "container image for applier to be used in the IMAGE:TAG format")
	applierImagePullPolicy := flag.String("image-pull-policy", "", "image pull policy for the applier image e.g. IfNotPresent, Always, etc.")
	sshKeySecret := flag.String("ssh-key-secret", "", "Kubernetes secret containing SSH key for the `id_rsa`")
	homeConfigMap := flag.String("home-configmap", "", "Kubernetes configmap containing applier's home-dir contents")
	helmX := flag.Bool("helm-x", true, "Enable helm-x integration for kustomize, k8s manifests, patching support")

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

	resource := schema.GroupVersionKind{
		Group:   "apps.mumoshu.github.io",
		Kind:    "Helmfile",
		Version: "v1alpha1",
	}

	if err := controller_runtime.Run(
		*name,
		resource,
		controller_runtime.Conf(*configPath),
		controller_runtime.Image(*applierImage),
		controller_runtime.ImagePullPolicy(*applierImagePullPolicy),
		controller_runtime.SSHKeySecret(*sshKeySecret),
		controller_runtime.HelmX(*helmX),
		controller_runtime.HomeConfigMap(*homeConfigMap),
	); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
