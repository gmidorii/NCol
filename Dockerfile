FROM golang:1.7-alpine

RUN echo "ipv6" >> /etc/modules
RUN apk update
RUN apk --no-cache add git
RUN apk --no-cache add make gcc g++

ENV SRCPATH /go/src/github.com/midorigreen/NCol
RUN mkdir -p $SRCPATH

COPY ./* $SRCPATH/

WORKDIR $SRCPATH
RUN make setup \
    && make deps \
		&& make update \
		&& make build

ENTRYPOINT ./NCol
