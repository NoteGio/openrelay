FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/websockets /websockets

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/websockets", "redis:6379", "topic://released", "postgres://postgres@postgres", "/run/secrets/postgress_password"]
