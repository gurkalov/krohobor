FROM golang:1.27.2-alpine3.24 as builder

WORKDIR /go/src/github.com/gurkalov/krohobor

COPY ./src/go.mod .
COPY ./src/go.sum .

RUN go mod download

COPY ./src .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/krohobor ./cmd


FROM ubuntu:26.04

RUN apt-get update
RUN apt-get install -y wget gnupg2

RUN apt-get install -y postgresql-common ca-certificates
RUN /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y

RUN apt-get -y update
RUN apt-get install -y postgresql-client-18 zip

WORKDIR /root

RUN mkdir /tmp/backup

COPY --from=builder /go/bin/krohobor .
RUN ln -s /root/krohobor /usr/bin/krohobor
