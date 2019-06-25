FROM golang:1.12 AS build

ADD . /helmfile-applier

WORKDIR /helmfile-applier

RUN GO111MODULES=on CGO_ENABLED=0 go build --mod=vendor -o helmfile-applier ./cmd

RUN find /helmfile-applier

RUN ./helmfile-applier -h

FROM quay.io/roboll/helmfile:v0.79.0

RUN helm init --client-only &&   \
  helm plugin install https://github.com/mumoshu/helm-x

# For private git repository access
RUN apk add --update openssh-client

COPY --from=build /helmfile-applier/helmfile-applier /usr/local/bin/

RUN find /usr/local/bin/

RUN /usr/local/bin/helmfile-applier -h

CMD /usr/local/bin/helmfile-applier
