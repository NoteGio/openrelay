FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/delayrelay /delayrelay

CMD ["/delayrelay"]
