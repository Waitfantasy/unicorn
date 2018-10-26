FROM golang:1.11.1-alpine3.8

RUN mkdir /etc/unicorn && mkdir /var/log/unicorn

WORKDIR /go/bin

COPY bin/unicorn .

CMD ["unicorn"]