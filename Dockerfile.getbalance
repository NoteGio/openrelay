FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/getbalance /getbalance

CMD ["/getbalance"]
