# This base image assumes that the bin directory has not been built previously
# and will build it for downstream consumption

FROM golang:1.8 as corebuild

RUN mkdir -p /go/src/github.com/notegio/openrelay

WORKDIR /go/src/github.com/notegio/openrelay

COPY . .

RUN make clean bin
