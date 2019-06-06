package controller

import (
	"fmt"
	"github.com/summerwind/whitebox-controller/config"
	"github.com/summerwind/whitebox-controller/handler"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewController(name string, client client.Client) (*config.ControllerConfig, error) {
	return &config.ControllerConfig{
		Name: name,
		Resource: schema.GroupVersionKind{
			Group: "apps.mumoshu.github.io",
			Kind: "Appliance",
			Version: "v1alpha1",
		},
		Dependents: []config.DependentConfig{
			{
				GroupVersionKind: schema.GroupVersionKind{
					Group:   "apps",
					Kind:    "Deployment",
					Version: "v1",
				},
				Orphan: false,
			},
		},
		Reconciler: &config.ReconcilerConfig{
			HandlerConfig: config.HandlerConfig{
				Func: &config.FuncHandlerConfig{
					Handler: NewReconcilingHandler(name, client),
				},
			},
		},
		Finalizer: &config.HandlerConfig{
			Func: &config.FuncHandlerConfig{
				Handler: NewFinalizingHandler(name, client),
			},
		},
	}, nil
}

type reconciclingHandler struct {
	name string
	c client.Client
}

func (h *reconciclingHandler) Run(buf []byte) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "reconciling: %s", string(buf))
	return nil, nil
}

type finalizingHandler struct {
	name string
	c client.Client
}

func (h *finalizingHandler) Run(buf []byte) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "finalizing: %s", string(buf))
	return nil, nil
}

func NewReconcilingHandler(name string, c client.Client) handler.Handler {
	return &reconciclingHandler{name: name, c: c}
}

func NewFinalizingHandler(name string, c client.Client) handler.Handler {
	return &finalizingHandler{name: name, c: c}
}
