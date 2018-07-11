FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/spendrecorder /spendrecorder

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/spendrecorder", "redis:6379", "queue://recordspend", "postgres://postgres@postgres", "/run/secrets/postgress_password"]
