package helmfile_operator

import (
	"github.com/go-logr/logr"
	"github.com/summerwind/whitebox-controller/config"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Option func(runner *GenericOperator) error

func Name(s string) Option {
	return func(r *GenericOperator) error {
		r.name = s
		return nil
	}
}

func Logger(l logr.Logger) Option {
	return func(r *GenericOperator) error {
		r.logger = l
		return nil
	}
}

func KubeConfig(c *rest.Config) Option {
	return func(r *GenericOperator) error {
		r.kc = c
		return nil
	}
}

func Manager(c manager.Manager) Option {
	return func(r *GenericOperator) error {
		r.mgr = c
		return nil
	}
}

func Config(c *config.Config) Option {
	return func(r *GenericOperator) error {
		r.c = c
		return nil
	}
}

func Controller(c *config.ControllerConfig) Option {
	return func(r *GenericOperator) error {
		r.controller = c
		return nil
	}
}

func Webhook(c *config.WebhookConfig) Option {
	return func(r *GenericOperator) error {
		r.webhook = c
		return nil
	}
}

