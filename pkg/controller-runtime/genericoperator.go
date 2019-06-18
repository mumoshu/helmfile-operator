package controller_runtime

import (
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/summerwind/whitebox-controller/config"
	"github.com/summerwind/whitebox-controller/controller"
	"github.com/summerwind/whitebox-controller/webhook"
)

type GenericOperator struct {
	name       string
	configPath string
	logger     logr.Logger

	webhook    *config.WebhookConfig
	controller *config.ControllerConfig
	kc         *rest.Config
	mgr        manager.Manager
	c          *config.Config
}

func New(opts ...Option) (*GenericOperator, error) {
	o := &GenericOperator{
	}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	return o, nil
}

func (o *GenericOperator) Run() error {
	c := o.c

	if o.controller != nil {
		if len(c.Controllers) == 0 {
			o.logger.Info("using build-in controller handler due to that nothing provided via config file", "config", o.configPath)
			c.Controllers = []*config.ControllerConfig{
				o.controller,
			}
		} else {
			o.logger.Info("using controllers provided via config file instead of built-in", "controllers", len(c.Controllers), "config", o.configPath)
		}
	}

	if o.webhook != nil {
		if c.Webhook == nil {
			o.logger.Info("using build-in webhook handler due to that nothing provided via config file", "config", o.configPath)
			c.Webhook = o.webhook
		} else {
			o.logger.Info("using webhook handler provided via config file instead of built-in", "config", o.configPath)
		}
	}

	for i, _ := range c.Controllers {
		cc := c.Controllers[i]
		_, err := controller.New(cc, o.mgr)
		if err != nil {
			return fmt.Errorf("could not create controller: %v", err)
		}
	}

	if c.Webhook != nil {
		_, err := webhook.NewServer(c.Webhook, o.mgr)
		if err != nil {
			return fmt.Errorf("could not create webhook server: %v", err)
		}
	}

	if err := o.mgr.Start(signals.SetupSignalHandler()); err != nil {
		return fmt.Errorf("could not start manager: %v", err)
	}

	return nil
}
