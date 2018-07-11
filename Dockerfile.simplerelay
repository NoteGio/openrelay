FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/simplerelay /simplerelay

CMD ["/simplerelay"]
