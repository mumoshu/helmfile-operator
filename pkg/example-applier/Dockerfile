FROM golang:1.12 AS build

ADD . /example-applier

WORKDIR /example-applier

RUN GO111MODULES=off go get -u github.com/gobuffalo/packr/v2/packr2

RUN GO111MODULES=on packr2

RUN GO111MODULES=on CGO_ENABLED=0 go build --mod=vendor -o example-applier ./cmd

RUN find /example-applier

RUN ./example-applier -h

FROM quay.io/roboll/helmfile:v0.73.1

COPY --from=build /example-applier/example-applier /usr/local/bin/

RUN find /usr/local/bin/

RUN /usr/local/bin/example-applier -h

CMD /usr/local/bin/example-applier
