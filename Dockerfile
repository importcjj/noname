FROM golang:1.18.0 AS builder

WORKDIR /myapp

COPY . .

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -o ddmc_bot ./cmd/maicai/main.go

FROM alpine:3.5

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /myapp
COPY --from=builder /myapp/ddmc_bot .
COPY --from=builder /myapp/sign.js .

CMD [ "./ddmc_bot", "-h" ]
