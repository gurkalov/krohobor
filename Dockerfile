FROM golang:1.13 as builder

WORKDIR /go/src/github.com/sr2020/krohobor

COPY ./src .

RUN GO111MODULE=on go get ./...

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/krohobor .


FROM alpine:latest

RUN apk --update add postgresql-client

WORKDIR /root/

COPY --from=builder /go/bin/krohobor .

EXPOSE 80

CMD ["./krohobor"]
