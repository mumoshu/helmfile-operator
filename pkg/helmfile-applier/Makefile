image:
	go mod vendor
	docker build -t mumoshu/helmfile-applier:dev .
	docker tag mumoshu/helmfile-applier:dev localhost:32000/helmfile-applier:dev
	docker push localhost:32000/helmfile-applier:dev

run:
	docker run --rm -it mumoshu/helmfile-applier:dev sh
