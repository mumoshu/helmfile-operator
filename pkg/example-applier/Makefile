image:
    packr2
	go mod vendor
	docker build -t mumoshu/example-applier:dev .
    packr2 clean

run:
	docker run --rm -it mumoshu/example-applier:dev sh
