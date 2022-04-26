FROM golang:1.18.0 AS builder

WORKDIR /myapp

COPY . .

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -o ddmc_bot ./cmd/maicai/main.go

FROM alpine:3.5

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

WORKDIR /myapp
COPY --from=builder /myapp/ddmc_bot .
COPY --from=builder /myapp/sign.js .
COPY --from=builder /myapp/config.yml .

CMD [ "./ddmc_bot", "-h" ]
