FROM golang:1.17.5-alpine AS builder
WORKDIR /go/src/github.com/xiazemin/registrator-nacos/
ENV GO111MODULE=on \
    GOPROXY="https://goproxy.cn"
COPY . .
RUN \
    CGO_ENABLED=0 go build \
		-a -installsuffix cgo \
		-ldflags "-X main.Version=$(cat VERSION)" \
		-o bin/registrator \
		.

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/xiazemin/registrator-nacos/bin/registrator /bin/registrator

ENTRYPOINT ["/bin/registrator"]
