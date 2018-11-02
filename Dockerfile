FROM golang:1.11.1-alpine3.8

RUN mkdir /etc/unicorn && mkdir /var/log/unicorn

RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /go/bin

COPY bin/unicorn .

ENV GIN_MODE release

EXPOSE 6001

EXPOSE 6002

CMD ["unicorn"]