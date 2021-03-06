FROM golang:alpine AS builder
RUN apk --no-cache add bash ca-certificates git make

# related to https://github.com/golang/go/issues/26988
ENV CGO_ENABLED 0

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega/...

COPY . /root/irc-notification-resource
RUN cd /root/irc-notification-resource && make test

FROM alpine:edge AS resource
COPY --from=builder /root/irc-notification-resource/artifacts /opt/resource
