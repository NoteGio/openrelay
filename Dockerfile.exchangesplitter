FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/exchangesplitter /exchangesplitter

CMD ["/exchangesplitter"]
