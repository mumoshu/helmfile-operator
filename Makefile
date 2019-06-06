# vgo needs to be enabled explicitly in packr2
# see https://github.com/gobuffalo/packr/issues/113
run: TARGET=examplecontroller
run:
	rm -rf dist/assets
	(GO111MODULE=on packr2 build github.com/mumoshu/appliance-operator/pkg/$(TARGET) && mv $(TARGET) ./dist && cd dist && ./$(TARGET))

example-applier/run:
	make run TARGET=example-applier

helmfile-applier/run:
	rm -rf dist/assets
	go build -o helmfile-applier ./pkg/helmfile-applier/cmd && mv helmfile-applier ./dist && cd dist && ./helmfile-applier

helmfile-operator/run:
	rm -rf dist/assets
	go build -o helmfile-operator ./pkg/helmfile-operator/cmd && mv helmfile-operator ./dist && cd dist && ./helmfile-operator

appliance-operator/run:
	rm -rf dist/assets
	go build -o appliance-operator ./pkg/appliance-operator && mv appliance-operator ./dist && cd dist && ./appliance-operator
