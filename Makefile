# vgo needs to be enabled explicitly in packr2
# see https://github.com/gobuffalo/packr/issues/113
run: TARGET=examplecontroller
run:
	rm -rf dist/assets
	(GO111MODULE=on packr2 build github.com/mumoshu/helmfile-server/pkg/$(TARGET) && mv $(TARGET) ./dist && cd dist && ./$(TARGET))

examplecontroller/run:
	make run TARGET=examplecontroller

genericcontroller/run:
	rm -rf dist/assets
	go build -o genericcontroller ./pkg/genericcontroller/cmd && mv genericcontroller ./dist && cd dist && ./genericcontroller

genericoperator/run:
	rm -rf dist/assets
	go build -o genericoperator ./pkg/genericoperator/cmd && mv genericoperator ./dist && cd dist && ./genericoperator

applianceoperator/run:
	rm -rf dist/assets
	go build -o applianceoperator ./pkg/applianceoperator && mv applianceoperator ./dist && cd dist && ./applianceoperator
