module github.com/mumoshu/helmfile-operator

go 1.12

require (
	github.com/davidovich/summon v0.7.0
	github.com/go-logr/logr v0.1.0
	github.com/gobuffalo/packr/v2 v2.3.1
	github.com/pkg/errors v0.8.1
	github.com/roboll/helmfile v0.73.0
	github.com/spf13/cobra v0.0.4
	github.com/stefanprodan/k8s-podinfo v1.4.2
	github.com/summerwind/whitebox-controller v0.5.0
	go.uber.org/zap v1.9.1
	k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v0.2.0
	sigs.k8s.io/controller-runtime v0.2.0-alpha.0
)

replace github.com/summerwind/whitebox-controller => ./whitebox-controller
