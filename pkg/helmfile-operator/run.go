package helmfile_operator

import (
	"fmt"
	"github.com/mumoshu/helmfile-server/pkg/helmfile-operator/controller"
	config2 "github.com/summerwind/whitebox-controller/config"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	log2 "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func Run(name string, o ...string) error {
	log2.SetLogger(log2.ZapLogger(false))
	log := log2.Log.WithName(name)

	var configPath string
	if len(o) > 0 {
		configPath = o[0]
	}
	kc, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("could not load kubernetes configuration: %v", err)
	}

	var c *config2.Config

	if configPath != "" {
		c, err = config2.LoadFile(configPath)
		if err != nil {
			return fmt.Errorf("could not load configuration file: %v", err)
		}
	} else {
		c = &config2.Config{}
	}

	opts := manager.Options{}
	if c.Metrics != nil {
		opts.MetricsBindAddress = fmt.Sprintf("%s:%d", c.Metrics.Host, c.Metrics.Port)
	}

	mgr, err := manager.New(kc, opts)
	if err != nil {
		return fmt.Errorf("could not create manager: %v", err)
	}

	controllerConfig, err := controller.NewController(name, mgr.GetClient())
	if err != nil {
		return fmt.Errorf("cloud not create controller config: %v", err)
	}

	operator, err := New(
		Config(c),
		Name(name),
		Logger(log),
		KubeConfig(kc),
		Manager(mgr),
		Controller(controllerConfig),
	)
	if err != nil {
		return fmt.Errorf("could not create operator: %v", err)
	}

	if err := operator.Run(); err != nil {
		return fmt.Errorf("could not run operator: %v", err)
	}

	return nil
}
