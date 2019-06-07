image:
	GO111MODULES=off go get -u github.com/gobuffalo/packr/v2/packr2
	packr2
	go mod vendor
	docker build -t mumoshu/helmfile-applier:dev .

run:
	docker run --rm -it mumoshu/helmfile-applier:dev sh
