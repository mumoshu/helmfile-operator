FROM golang:1.12 AS build

ADD . /helmfile-applier

WORKDIR /helmfile-applier

RUN GO111MODULES=on CGO_ENABLED=1 go build --mod=vendor -o helmfile-applier ./cmd

RUN find /helmfile-applier

RUN ./helmfile-applier -v

FROM quay.io/roboll/helmfile:v0.73.1

COPY --from=build /helmfile-applier/helmfile-applier /usr/local/bin/

RUN find /usr/local/bin/

CMD /usr/local/bin/helmfile-applier
