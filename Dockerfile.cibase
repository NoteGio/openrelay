# This base image assumes that the bin directory will have been built
# previously and can be copied into the image

FROM golang:1.8 as corebuild

RUN mkdir -p /go/src/github.com/notegio/openrelay

COPY ./bin /go/src/github.com/notegio/openrelay/bin
