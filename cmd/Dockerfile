# Build ltk in a stock Go builder
FROM golang:1.13-alpine as builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache make musl-dev linux-headers git

ENV SRC_DIR=/ltk
COPY . $SRC_DIR
WORKDIR $SRC_DIR

RUN mkdir -p /go_bin
RUN export GO111MODULE=on && export GOPROXY=https://goproxy.cn && go mod tidy && \
    go generate $SRC_DIR/pkg/version && \
    CGO_ENABLED=1 GOBIN=/go_bin go install -ldflags '-w -s' -v -tags netgo $SRC_DIR/cmd/...

# Pull ltk into a second stage deploy alpine container
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates

COPY --from=builder /go_bin/* /usr/local/bin/

CMD ["sh"]