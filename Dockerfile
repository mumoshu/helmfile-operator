FROM golang:1.12.4-alpine3.9 as builder

RUN apk add --no-cache make git

RUN GO111MODULES=off go get -u github.com/gobuffalo/packr/v2/packr2

WORKDIR /workspace/helmfile-operator

ADD ./ /workspace/helmfile-operator
RUN make indocker-build

# -----------------------------------------------------------------------------

FROM quay.io/roboll/helmfile:v0.79.0

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.14.2/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl

COPY --from=builder /workspace/helmfile-operator/helmfile-operator /usr/local/bin/helmfile-operator
