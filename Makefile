deps:
	GO111MODULES=off go get -u github.com/gobuffalo/packr/v2/packr2


# vgo needs to be enabled explicitly in packr2
# see https://github.com/gobuffalo/packr/issues/113
run: TARGET=examplecontroller
run:
	rm -rf dist/assets
	(GO111MODULES=on packr2 build github.com/mumoshu/helmfile-operator/pkg/$(TARGET) && mv $(TARGET) ./dist && cd dist && ./$(TARGET))

example-applier/run:
	make run TARGET=example-applier

helmfile-applier/run:
	rm -rf dist/assets
	go build -o helmfile-applier ./pkg/helmfile-applier/cmd && mv helmfile-applier ./dist && cd dist && ./helmfile-applier

controller-runtime/run:
	rm -rf dist/assets
	go build -o controller-runtime ./pkg/controller-runtime/cmd && mv controller-runtime ./dist && cd dist && ./controller-runtime

helmfile-operator/run:
	rm -rf dist/assets
	go build -o helmfile-operator ./pkg/helmfile-operator && mv helmfile-operator ./dist && cd dist && ./helmfile-operator

build:
	go build -o helmfile-operator ./pkg/helmfile-operator
	go build -o controller-runtime ./pkg/controller-runtime/cmd
	go build -o helmfile-applier ./pkg/helmfile-applier/cmd
