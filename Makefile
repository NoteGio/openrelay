PACKAGE  = github.com/notegio/openrelay
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)
GOSTATIC = go build -a -installsuffix cgo -ldflags '-extldflags "-static"'

all: bin nodesetup truffleCompile docker-cfg/ca-certificates.crt

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

clean: dockerstop
	rm -rf bin/ .gopath/
	rm -rf js/build


dockerstop:
	docker stop `cat "$(BASE)/tmp/redis.containerid"` || true
	docker rm `cat "$(BASE)/tmp/redis.containerid"` || true
	rm $(BASE)/tmp/redis.containerid || true
	docker stop `cat "$(BASE)/tmp/postgres.containerid"` || true
	docker rm `cat "$(BASE)/tmp/postgres.containerid"` || true
	rm $(BASE)/tmp/postgres.containerid || true

nodesetup:
	cd js ; npm install

bin/delayrelay: $(BASE) cmd/delayrelay/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/delayrelay cmd/delayrelay/main.go

bin/fundcheckrelay: $(BASE) cmd/fundcheckrelay/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/fundcheckrelay cmd/fundcheckrelay/main.go

bin/getbalance: $(BASE) cmd/getbalance/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/getbalance cmd/getbalance/main.go

bin/ingest: $(BASE) cmd/ingest/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/ingest cmd/ingest/main.go

bin/initialize: $(BASE) cmd/initialize/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/initialize cmd/initialize/main.go

bin/simplerelay: $(BASE) cmd/simplerelay/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/simplerelay cmd/simplerelay/main.go

bin/validateorder: $(BASE) cmd/validateorder/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/validateorder cmd/validateorder/main.go

bin/fillupdate: $(BASE) cmd/fillupdate/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/fillupdate cmd/fillupdate/main.go

bin/indexer: $(BASE) cmd/indexer/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/indexer cmd/indexer/main.go

bin/fillindexer: $(BASE) cmd/fillindexer/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/fillindexer cmd/fillindexer/main.go

bin/blockmonitor: $(BASE) cmd/blockmonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/blockmonitor cmd/blockmonitor/main.go

bin/allowancemonitor: $(BASE) cmd/allowancemonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/allowancemonitor cmd/allowancemonitor/main.go

bin/erc721approvalmonitor: $(BASE) cmd/erc721approvalmonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/erc721approvalmonitor cmd/erc721approvalmonitor/main.go

bin/canceluptomonitor: $(BASE) cmd/canceluptomonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/canceluptomonitor cmd/canceluptomonitor/main.go

bin/canceluptofilter: $(BASE) cmd/canceluptofilter/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/canceluptofilter cmd/canceluptofilter/main.go

bin/canceluptoindexer: $(BASE) cmd/canceluptoindexer/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/canceluptoindexer cmd/canceluptoindexer/main.go

bin/affiliatemonitor: $(BASE) cmd/affiliatemonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/affiliatemonitor cmd/affiliatemonitor/main.go

bin/spendmonitor: $(BASE) cmd/spendmonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/spendmonitor cmd/spendmonitor/main.go

bin/fillmonitor: $(BASE) cmd/fillmonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/fillmonitor cmd/fillmonitor/main.go

bin/multisigmonitor: $(BASE) cmd/multisigmonitor/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/multisigmonitor cmd/multisigmonitor/main.go

bin/spendrecorder: $(BASE) cmd/spendrecorder/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/spendrecorder cmd/spendrecorder/main.go

bin/exchangesplitter: $(BASE) cmd/exchangesplitter/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/exchangesplitter cmd/exchangesplitter/main.go

bin/automigrate: $(BASE) cmd/automigrate/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/automigrate cmd/automigrate/main.go

bin/searchapi: $(BASE) cmd/searchapi/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/searchapi cmd/searchapi/main.go

bin/queuemonitor: $(BASE) cmd/queuemonitor/main.go
	cd "$(BASE)" && CGO_ENABLED=0 $(GOSTATIC) -o bin/queuemonitor cmd/queuemonitor/main.go

bin/terms: $(BASE) cmd/terms/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/terms cmd/terms/main.go

bin/poolfilter: $(BASE) cmd/poolfilter/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/poolfilter cmd/poolfilter/main.go

bin/metadataindexer: $(BASE) cmd/metadataindexer/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/metadataindexer cmd/metadataindexer/main.go

bin/websockets: $(BASE) cmd/websockets/main.go
	cd "$(BASE)" && $(GOSTATIC) -o bin/websockets cmd/websockets/main.go

bin: bin/delayrelay bin/fundcheckrelay bin/getbalance bin/ingest bin/initialize bin/simplerelay bin/validateorder bin/fillupdate bin/indexer bin/fillindexer bin/automigrate bin/searchapi bin/exchangesplitter bin/blockmonitor bin/allowancemonitor bin/spendmonitor bin/fillmonitor bin/multisigmonitor bin/spendrecorder bin/queuemonitor bin/canceluptomonitor bin/canceluptofilter bin/canceluptoindexer bin/erc721approvalmonitor bin/affiliatemonitor bin/terms bin/poolfilter bin/metadataindexer bin/websockets

truffleCompile:
	cd js ; node_modules/.bin/truffle compile

$(BASE)/tmp/redis.containerid:
	mkdir -p $(BASE)/tmp
	docker run -d -p 6379:6379 redis  > $(BASE)/tmp/redis.containerid

$(BASE)/tmp/postgres.containerid:
	mkdir -p $(BASE)/tmp
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=secret postgres > $(BASE)/tmp/postgres.containerid

dockerstart: $(BASE) $(BASE)/tmp/redis.containerid $(BASE)/tmp/postgres.containerid

gotest: dockerstart test-funds test-channels test-accounts test-affiliates test-types test-ingest test-blocksmonitor test-allowancemonitor test-fillmonitor test-spendmonitor test-splitter test-search test-db test-metadata test-pool test-ws test-subscriptions

test-funds: $(BASE)
	cd "$(BASE)/funds" && go test
test-channels: $(BASE)
	cd "$(BASE)/channels" &&  REDIS_URL=localhost:6379 go test
test-accounts: $(BASE)
	cd "$(BASE)/accounts" &&  REDIS_URL=localhost:6379 go test
test-affiliates: $(BASE)
	cd "$(BASE)/affiliates" &&  REDIS_URL=localhost:6379 go test
test-types: $(BASE)
	cd "$(BASE)/types" && go test
test-ingest: $(BASE)
	cd "$(BASE)/ingest" && go test
test-blocksmonitor: $(BASE)
	cd "$(BASE)/monitor/blocks" && go test
test-allowancemonitor: $(BASE)
	cd "$(BASE)/monitor/allowance" && go test
test-erc721approval: $(BASE)
	cd "$(BASE)/monitor/erc721approval" && go test
test-canceluptomonitor: $(BASE)
	cd "$(BASE)/monitor/cancelupto" && go test
test-fillmonitor: $(BASE)
	cd "$(BASE)/monitor/fill" && go test
test-spendmonitor: $(BASE)
	cd "$(BASE)/monitor/spend" && go test
test-splitter: $(BASE)
	cd "$(BASE)/splitter" && go test
test-search: $(BASE)
	cd "$(BASE)/search" && POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=secret go test
test-affiliate: $(BASE)
	cd "$(BASE)/monitor/affiliate" && go test
test-db: $(BASE)
	cd "$(BASE)/db" &&  POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=secret go test
test-metadata: $(BASE)
	cd "$(BASE)/metadata" &&  POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=secret go test
test-pool: $(BASE)
	cd "$(BASE)/pool" &&  POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=secret go test
test-ws: $(BASE)
	cd "$(BASE)/channels/ws" &&  POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=secret go test
test-subscriptions: $(BASE)
	cd "$(BASE)/subscriptions" &&  POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=secret go test

docker-cfg/ca-certificates.crt:
	cp /etc/ssl/certs/ca-certificates.crt docker-cfg/ca-certificates.crt

test: $(BASE) $(BASE)/tmp/redis.containerid gotest dockerstop
test_no_docker: mock gotest
mock: $(BASE)
	mkdir -p $(BASE)/tmp
	touch $(BASE)/tmp/redis.containerid
	touch $(BASE)/tmp/postgres.containerid
newvendor:
	govendor add +external

0x-testrpc-snapshot.tar.gz:
	wget https://s3.amazonaws.com/testrpc-shapshots/965d6098294beb22292090c461151274ee6f9a26.zip -O testrpc-db.zip
	mkdir -p /tmp/testrpc-snapshot
	unzip testrpc-db.zip -d /tmp/testrpc-snapshot
	tar -czf 0x-testrpc-snapshot.tar.gz -C /tmp/testrpc-snapshot/0x_ganache_snapshot .
	rm testrpc-db.zip
	rm -rf /tmp/testrpc-snapshot
