# appliance-operator

`appliance-operator` is a Kubernetse operator that deploys `Appliance`s to your K8s cluster.

The benefit of `Appliance` is that you can declaratively deploy K8s manifests generated with any tool from any source:

Supported tools:

- Vanilla Kubernetes manifests(Regular Kuberntees `YAML`s)
- Helm
- Kustomize
- Helm + Kustomize(JIT kustomize patches before installinh Helm charts)

Supported sources:

- Git
- AWS S3
- Potentially any sources supported by [go-getter](https://github.com/hashicorp/go-getter) and helm downloader plugins

## The "Appliance" custom resource

Seeing is believing - here's the custom resource definition of `Appliance` for you to get what it makes possible.

It will be a lot easier if you read it along with [the official documentation for CustomResourceDefinition](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#create-a-customresourcedefinition):

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: appliance.apps.mumoshu.github.io
spec:
  group: apps.mumoshu.github.io
  versions:
  - name: v1alpha1
    served: true
    storage: true
  names:
    kind: Appliance
    plural: appliance
    singular: appliance
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        spec:
          properties:
            source:
              type: string
            values:
              type: object
            envvars:
              type: object
```

And an example custom resource:

```yaml
apiVersion: apps.mumoshu.github.io/v1alpha1
kind: Appliance
metadata:
  name: myapp
spec:
  source: git::https://github.com/mumoshu/appliance-operator//pkg/examplecontroller@assets?ref=master
  #
  # # Syntax sugar: The trailing "?ref=master" can be specified alternatively with `version`:
  # version: master
  #
  # # Alternatively specify the container `image` for your own appliance-controller.
  # # See the README section for `examplecontroller` to see how you can build your own appliance-controller.
  # image:
  #   repository: quay.io/examplecom/example
  #   tag: v1.2.3
  #
  values:
    foo: FOO
  envvars:
    bar: BAR
```

## Comparison to similar systems

Is `appliance-operator` the only solution to the problem stated in the top of this page?

No - but I believe only `appliance-operator` does it almost nearly perfect :)

Let's compare it with existing and similar systems, so that you can get why and in which circumstance `appliance-operator`(or other system) is good fit:

- [addons-operators](https://github.com/kubernetes-sigs/addon-operators) strives to be the collection of operators for `cluster addons` while providing the solid framework for building those operators.

  `appliance-operator`, on the other hand, is for managing `appliance`. An appliance is not necessarily a cluster addon so it was natural for the author of this project(me) to build it outside the cluster-operators project.
  
  An cluster-addon operator is free to use `appliance-operator` and delegate the management of a complex of bundle of K8s maniefsts(=an `appliance`) to `appliance-operator`.
  
  The interface between the addon operator and the appliance-operator will be `Appliance` custom resources. The addon operator should CRUD `Appliance` resources via K8s API accordingly to the `Addon` and `Channel` resources.
  So that each `addon-operator` can be maintained easily without dealing with various K8s deployment tools.    

## Terminology

`appliance` is a virtual appliance a.k.a a self-managed service that is installed onto your k8s cluster.
An `appliance` is composed of one or more K8s resources.
 
`consumer` is the user of the appliance. An `appliance` surface a single U/X to manage the appliance as a whole, so that `consumer` is able to operate on the `appliance`, rather than on underlying K8s resources.

`applier` periodically reconcile the k8s cluster so that the appliance keeps running.
This project aims to provide (1) generic `applier` and (2) the framework for building your own `applier` and `appliance-operator`.

`helmfile-applier` is the only `applier` that is maintained in this project. `example-applier` is an example applier that shows how to build your own applier based on `helmfile-applier`.

`appliance-operator` watches for `Appliance` custom resources so that it can reconcile the k8s cluster to have corresponding `appliance-controller` up and running.

## Projects

There are 1 main project and 2 related projects maintained within this project.

See respective directory under the `pkg` directory for more details.

### `appliance-operator`

This is the example impelmentation of your own `appliance-operator`.

It is based on the `helmfile-operator`, and has the additional ability to automatically register the `Appliance` CRD on startup.

It's configuration can be customized via the changing the bundled `assets/config.yaml`, or mouting an alternative config file and pointing the operator to load that by providing the `--file CONFIG_FILE` flag.

For settings available in `config.yaml`, see the documentation of [whitebox-controller](https://github.com/summerwind/whitebox-controller/blob/master/docs/configuration.md) which is the foundation of the operator.

### `helmfile-operator`

This is a generic implementation of the `appliance-operator` which is deployed onto a runtime environment that has connectivity to K8s API.

By contacting K8s API, it watches `Appliance` custom resources for changes and updates and removes the `Deployment` that runs the `helmfile-applier`.

See [The "Appliance" custom resource](#The "Appliance" custom resource) for details of the `Appliance` resource.

### helmfile-applier

This is a generic `applier` that is deployed as a container image.

You add the state file for the appliance into the container image so that `helmfile-applier` loads it on startup. After that it periodically runs [helmfile](https://github.com/roboll/helmfile) to reconcile the cluster state so that
the desired appliance is kept up and running.

### example-applier

Use-cases: Air-gapped deployments (By containerizing all the appliance assets)

This is an example `applier` that is implemented by extending `helmfile-applier`.

It bundles the state file for the appliance into the executable binary with [packr2](https://github.com/gobuffalo/packr).

On runtime, it calls [summon](https://github.com/davidovich/summon) to extract the state file and its belongings, periodically run [helmfile](https://github.com/roboll/helmfile) commands to reconcile the cluster state so that
the desired appliance is kept up and running.

## Advanced: Building your own Appliance operator

TL;DR; Use `helmfile-operator` as a framework and `appliance-operator` as the reference implementation for building your own operator.

> **For advanced-usecases**
> 
> Consider this as a framework for building your own appliance-operator-like operator from scrach.

## Related projects

Under the hood, `appliance-operator` uses: 

- [helmfile](https://github.com/roboll/helmfile) for declaratively manage and deploy appliances onto your cluster.
- [helm-x](https://github.com/mumoshu/helm-x) for transparent support for any K8s manifests renderer/builder.
- [kustomize](https://github.com/kubernetes-sigs/kustomize) for JIT patching your K8s manifests or helm charts contained in the appliance.
- [helm](https://github.com/helm/helm/) for installing, upgrading your appliance
- [helm-diff](https://github.com/databus23/helm-diff) for reviewing changes before actually updating the cluster state.
- [whitebox-controller](https://github.com/summerwind/whitebox-controller) as the framework for building the rock-solid appliance-operator that is customizable via config files and shell scripts.
