FROM golang:1.17.2-alpine3.14 as builder

WORKDIR /go/src/github.com/gurkalov/krohobor

COPY ./src/go.mod .
COPY ./src/go.sum .

RUN go mod download

COPY ./src .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/krohobor ./cmd


FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install -y wget gnupg2

RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ focal-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update

RUN apt-get install -y postgresql-client-14 zip

WORKDIR /root

RUN mkdir /tmp/backup

COPY --from=builder /go/bin/krohobor .
RUN ln -s /root/krohobor /usr/bin/krohobor
