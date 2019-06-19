module github.com/mumoshu/helmfile-operator

go 1.12

require (
	cloud.google.com/go v0.40.0 // indirect
	github.com/Masterminds/sprig v2.20.0+incompatible // indirect
	github.com/aokoli/goutils v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.20.3 // indirect
	github.com/cheggaaa/pb v2.0.6+incompatible // indirect
	github.com/davidovich/summon v0.7.0
	github.com/go-logr/logr v0.1.0
	github.com/gobuffalo/packr/v2 v2.3.2
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/kr/pty v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mumoshu/helmfile-operator/pkg/helmfile-applier v0.0.0-20190618024137-7a3a8b54885f // indirect
	github.com/pkg/errors v0.8.1
	github.com/roboll/helmfile v0.79.0
	github.com/spf13/cobra v0.0.4
	github.com/stefanprodan/k8s-podinfo v1.4.2
	github.com/summerwind/whitebox-controller v0.5.0
	github.com/ulikunitz/xz v0.5.6 // indirect
	github.com/urfave/cli v1.20.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190618222545-ea8f1a30c443 // indirect
	golang.org/x/image v0.0.0-20190618124811-92942e4437e2 // indirect
	golang.org/x/mobile v0.0.0-20190607214518-6fa95d984e88 // indirect
	golang.org/x/net v0.0.0-20190619014844-b5b0513f8c1b // indirect
	golang.org/x/sys v0.0.0-20190618155005-516e3c20635f // indirect
	golang.org/x/tools v0.0.0-20190618233249-04b924abaa25 // indirect
	google.golang.org/genproto v0.0.0-20190611190212-a7e196e89fd3 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28 // indirect
	honnef.co/go/tools v0.0.0-20190614002413-cb51c254f01b // indirect
	k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v0.2.0
	sigs.k8s.io/controller-runtime v0.2.0-alpha.0
)

replace github.com/summerwind/whitebox-controller => ./whitebox-controller

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
