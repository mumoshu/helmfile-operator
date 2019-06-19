package controller_runtime

import (
	"fmt"
	"github.com/mumoshu/helmfile-operator/pkg/controller-runtime/controller"
	whiteboxctlr "github.com/summerwind/whitebox-controller/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	log2 "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type Runtime struct {
	source, image, imagePullPolicy, config, sshKeySecret string

	helmX bool
}

type Opt func(r *Runtime) error

func Source(s string) Opt {
	return func(r *Runtime) error {
		r.source = s
		return nil
	}
}

func Image(s string) Opt {
	return func(r *Runtime) error {
		r.image = s
		return nil
	}
}

func ImagePullPolicy(s string) Opt {
	return func(r *Runtime) error {
		r.imagePullPolicy = s
		return nil
	}
}

func Conf(s string) Opt {
	return func(r *Runtime) error {
		r.config = s
		return nil
	}
}

func SSHKeySecret(s string) Opt {
	return func(r *Runtime) error {
		r.sshKeySecret = s
		return nil
	}
}

func HelmX(b bool) Opt {
	return func(r *Runtime) error {
		r.helmX = b
		return nil
	}
}

func Run(name string, resource schema.GroupVersionKind, opt ...Opt) error {
	r := &Runtime{}
	for _, o := range opt {
		if err := o(r); err != nil {
			return err
		}
	}

	log2.SetLogger(log2.ZapLogger(false))
	log := log2.Log.WithName(name)

	var configPath string
	if r.config != "" {
		configPath = r.config
	}
	kc, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("could not load kubernetes configuration: %v", err)
	}

	var c *whiteboxctlr.Config

	if configPath != "" {
		c, err = whiteboxctlr.LoadFile(configPath)
		if err != nil {
			return fmt.Errorf("could not load configuration file: %v", err)
		}
	} else {
		c = &whiteboxctlr.Config{}
	}

	opts := manager.Options{}
	if c.Metrics != nil {
		opts.MetricsBindAddress = fmt.Sprintf("%s:%d", c.Metrics.Host, c.Metrics.Port)
	}

	mgr, err := manager.New(kc, opts)
	if err != nil {
		return fmt.Errorf("could not create manager: %v", err)
	}

	controllerConfig, err := controller.NewController(name, resource, mgr.GetClient(), r.source, r.image, r.imagePullPolicy, r.sshKeySecret, r.helmX)
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
