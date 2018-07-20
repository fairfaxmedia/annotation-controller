FROM golang:1.9.2-alpine3.6 as builder

ENV GOOS=linux GOARCH=amd64

WORKDIR /go/src/github.com/fairfaxmedia/annotation-controller/

COPY src/ /go/src/github.com/fairfaxmedia/annotation-controller/src/
COPY hack/ /go/src/github.com/fairfaxmedia/annotation-controller/hack/
COPY Gopkg.* /go/src/github.com/fairfaxmedia/annotation-controller/
COPY Makefile /go/src/github.com/fairfaxmedia/annotation-controller/
RUN apk add --no-cache \
        ca-certificates tzdata git curl bash make && \
        rm -rf /var/cache/apk/*

RUN curl -s -o /usr/local/bin/dep -L https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod 755 /usr/local/bin/dep

RUN make depend build-go

FROM alpine:3.6
RUN apk add --no-cache \
        ca-certificates tzdata && \
        rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/fairfaxmedia/annotation-controller/src/cmd/controller/controller /go/bin/
CMD cd /go/bin/ && ./controller
