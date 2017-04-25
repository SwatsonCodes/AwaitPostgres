FROM golang:1

ARG PACKAGE=github.com/swatsoncodes/AwaitPostgres/

ADD main.go /go/src/${PACKAGE}
WORKDIR /go/src/${PACKAGE}

RUN set -eux &&\
  go get -v -d . &&\
  go build -v -o /AwaitPostgres

CMD ["/AwaitPostgres"]
