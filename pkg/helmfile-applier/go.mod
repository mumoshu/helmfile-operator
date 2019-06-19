module github.com/mumoshu/helmfile-operator/pkg/helmfile-applier

go 1.12

require (
	github.com/gobuffalo/packr/v2 v2.3.2
	github.com/mumoshu/helmfile-operator v0.0.0-20190618020232-749c7de4d3a3
	github.com/roboll/helmfile v0.79.0
	github.com/spf13/cobra v0.0.4
	github.com/stefanprodan/k8s-podinfo v1.4.2
	go.uber.org/zap v1.10.0
)

replace github.com/mumoshu/helmfile-operator => ../../
