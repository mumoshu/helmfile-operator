package controller

import (
	"encoding/json"
	"fmt"
	"github.com/summerwind/whitebox-controller/config"
	"github.com/summerwind/whitebox-controller/handler"
	"github.com/summerwind/whitebox-controller/reconciler"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewController(name string, resource schema.GroupVersionKind, client client.Client, source, image, imagePullPolicy, sshKeySecret string) (*config.ControllerConfig, error) {
	if imagePullPolicy == "" {
		imagePullPolicy = "IfNotPresent"
	}

	h := &reconciclingHandler{
		name:            name,
		c:               client,
		source:          source,
		image:           image,
		imagePullPolicy: imagePullPolicy,
		sshKeySecret:    sshKeySecret,
	}

	return &config.ControllerConfig{
		Name:     name,
		Resource: resource,
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
					Handler: h,
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
	name            string
	c               client.Client
	source, image   string
	imagePullPolicy string
	sshKeySecret    string
}

const DefaultImageTag = "mumoshu/helmfile-applier:dev"

func (h *reconciclingHandler) Run(buf []byte) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "reconciling: %s\n", string(buf))

	state := reconciler.State{}

	if err := json.Unmarshal(buf, &state); err != nil {
		return buf, err
	}

	deployName := fmt.Sprintf("%s-%s", h.name, state.Object.Object["metadata"].(map[string]interface{})["name"])

	objectSpec := state.Object.Object["spec"].(map[string]interface{})
	dependentDeploys := state.Dependents["deployment.v1.apps"]

	command := "helmfile-applier"
	if v, ok := objectSpec["command"]; ok {
		cmd := v.(string)
		if cmd != "" {
			command = cmd
		}
	}

	args := []string{
		command,
	}

	env := []map[string]interface{}{}

	var source string
	if v, ok := objectSpec["source"]; ok {
		source = v.(string)
	} else if h.source != "" {
		source = h.source
	}
	if source != "" {
		args = append(args, "--file", source)
	}

	if v, ok := objectSpec["values"]; ok {
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		args = append(args, "--values", string(bytes))
	}

	if v, ok := objectSpec["valuesFiles"]; ok {
		switch t := v.(type) {
		case string:
			args = append(args, "--valuesFiles", t)
		case interface{}:
			args = append(args, "--valuesFiles", t.(string))
		default:
			return nil, fmt.Errorf("unexpected type of valuesFiles entry: value=%v, type=%T", t, t)
		}
	}

	var environment string
	if v, ok := objectSpec["environment"]; ok {
		environment = v.(string)

		args = append(args, "--environment", environment)
	}

	var envvars map[string]interface{}
	if v, ok := objectSpec["envvars"]; ok {
		envvars = v.(map[string]interface{})
		for name, val := range envvars {
			env = append(env, map[string]interface{}{
				"name":  name,
				"value": val,
			})
		}
	}

	volumes := []map[string]interface{}{}
	volumeMounts := []map[string]interface{}{}

	if h.sshKeySecret != "" {
		mount := map[string]interface{}{
			// TODO
			"name":      "dot-ssh",
			"mountPath": "/root/.ssh",
		}
		volume := map[string]interface{}{
			"name": "dot-ssh",
			"secret": map[string]interface{}{
				"secretName":  h.sshKeySecret,
				"defaultMode": 500,
			},
		}
		volumes = append(volumes, volume)
		volumeMounts = append(volumeMounts, mount)
	}

	if dependentDeploys == nil || len(dependentDeploys) == 0 {
		var imageTag string

		if _, ok := objectSpec["image"]; ok {
			image := objectSpec["image"].(map[string]interface{})
			repo := image["repository"].(string)
			tag := image["tag"].(string)
			imageTag = fmt.Sprintf("%s:%s", repo, tag)
		} else if h.image != "" {
			imageTag = h.image
		} else {
			imageTag = DefaultImageTag
		}

		// TODO Use https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/blob/master/pkg/patterns/declarative/pkg/manifest/objects.go
		state.Dependents["deployment.v1.apps"] = []*unstructured.Unstructured{
			&unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      deployName,
						"namespace": "default",
					},
					"spec": map[string]interface{}{
						"replicas": 1,
						"selector": map[string]interface{}{
							"matchLabels": map[string]interface{}{
								"app": deployName,
							},
						},
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"app": deployName,
								},
							},
							"spec": map[string]interface{}{
								"containers": []map[string]interface{}{
									map[string]interface{}{
										"name":            "helmfile-applier",
										"command":         args,
										"image":           imageTag,
										"imagePullPolicy": h.imagePullPolicy,
										"env":             env,
										"volumeMounts":    volumeMounts,
									},
								},
								"volumes": volumes,
							},
						},
					},
				},
			},
		}
	}

	out, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "new state: %s\n", string(out))

	return out, nil
}

type finalizingHandler struct {
	name string
	c    client.Client
}

func (h *finalizingHandler) Run(buf []byte) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "finalizing: %s\n", string(buf))

	state := reconciler.State{}

	if err := json.Unmarshal(buf, &state); err != nil {
		return buf, err
	}

	state.Dependents["deployment.v1.apps"] = []*unstructured.Unstructured{}

	out, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "new state: %s\n", string(out))

	return out, nil
}

func NewFinalizingHandler(name string, c client.Client) handler.Handler {
	return &finalizingHandler{name: name, c: c}
}
