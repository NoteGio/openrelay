FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/queuemonitor /queuemonitor

CMD ["/queuemonitor", "redis:6379", "1", "newblocks-ropsten"]
