PACKAGE  = github.com/notegio/openrelay
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)
GOSTATIC = go build -a -installsuffix cgo -ldflags '-extldflags "-static"'


all: bin nodesetup truffleCompile certs

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

clean: dockerstop
	rm -rf bin/ .gopath/
	rm -rf js/build


dockerstop:
	docker stop `cat $(BASE)/tmp/dynamo.containerid` || true
	docker rm `cat $(BASE)/tmp/dynamo.containerid` || true
	rm $(BASE)/tmp/dynamo.containerid || true
	docker stop `cat $(BASE)/tmp/redis.containerid` || true
	docker rm `cat $(BASE)/tmp/redis.containerid` || true
	rm $(BASE)/tmp/redis.containerid || true
	docker stop `cat $(BASE)/tmp/postgres.containerid` || true
	docker rm `cat $(BASE)/tmp/postgres.containerid` || true
	rm $(BASE)/tmp/postgres.containerid || true

nodesetup:
	cd js ; npm install

bin/delayrelay: $(BASE) cmd/delayrelay/main.go
	cd $(BASE) &&  CGO_ENABLED=0 $(GOSTATIC) -o bin/delayrelay cmd/delayrelay/main.go

bin/fundcheckrelay: $(BASE) cmd/fundcheckrelay/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/fundcheckrelay cmd/fundcheckrelay/main.go

bin/getbalance: $(BASE) cmd/getbalance/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/getbalance cmd/getbalance/main.go

bin/ingest: $(BASE) cmd/ingest/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/ingest cmd/ingest/main.go

bin/initialize: $(BASE) cmd/initialize/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/initialize cmd/initialize/main.go

bin/simplerelay: $(BASE) cmd/simplerelay/main.go
	cd $(BASE) && CGO_ENABLED=0 $(GOSTATIC) -o bin/simplerelay cmd/simplerelay/main.go

bin/validateorder: $(BASE) cmd/validateorder/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/validateorder cmd/validateorder/main.go

bin/fillupdate: $(BASE) cmd/fillupdate/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/fillupdate cmd/fillupdate/main.go

bin/indexer: $(BASE) cmd/indexer/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/indexer cmd/indexer/main.go

bin/fillindexer: $(BASE) cmd/fillindexer/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/fillindexer cmd/fillindexer/main.go

bin/automigrate: $(BASE) cmd/automigrate/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/automigrate cmd/automigrate/main.go

bin/searchapi: $(BASE) cmd/searchapi/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/searchapi cmd/searchapi/main.go

bin: bin/delayrelay bin/fundcheckrelay bin/getbalance bin/ingest bin/initialize bin/simplerelay bin/validateorder bin/fillupdate bin/indexer bin/fillindexer bin/automigrate bin/searchapi

truffleCompile:
	cd js ; node_modules/.bin/truffle compile

$(BASE)/tmp/redis.containerid:
	mkdir -p $(BASE)/tmp
	docker run -d -p 6379:6379 redis  > $(BASE)/tmp/redis.containerid

$(BASE)/tmp/postgres.containerid:
	mkdir -p $(BASE)/tmp
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=secret postgres > $(BASE)/tmp/postgres.containerid

$(BASE)/tmp/dynamo.containerid:
	mkdir -p $(BASE)/tmp
	docker run -d -p 8000:8000 cnadiminti/dynamodb-local > $(BASE)/tmp/dynamo.containerid

py/.env:
	virtualenv -p python3.6 $(BASE)/py/.env
	$(BASE)/py/.env/bin/pip install -r $(BASE)/py/requirements/api.txt
	$(BASE)/py/.env/bin/pip install -r $(BASE)/py/requirements/indexer.txt
	$(BASE)/py/.env/bin/pip install nose

gotest: $(BASE)/tmp/redis.containerid $(BASE)/tmp/postgres.containerid
	cd $(BASE)/funds && go test
	cd $(BASE)/channels &&  REDIS_URL=localhost:6379 go test
	cd $(BASE)/accounts &&  REDIS_URL=localhost:6379 go test
	cd $(BASE)/affiliates &&  REDIS_URL=localhost:6379 go test
	cd $(BASE)/types && go test
	cd $(BASE)/ingest && go test
	cd $(BASE)/search && go test
	cd $(BASE)/db &&  POSTGRES_HOST=localhost POSTGRES_USER=postgres POSTGRES_PASSWORD=secret go test

pytest: $(BASE)/tmp/dynamo.containerid
	cd $(BASE)/py && DYNAMODB_HOST="http://localhost:8000" $(BASE)/py/.env/bin/nosetests

jstest: $(BASE)/tmp/redis.containerid
	cd $(BASE)/js && REDIS_URL=localhost:6379 node_modules/.bin/mocha

certs: docker-cfg/ca-certificates.crt
	cp /etc/ssl/certs/ca-certificates.crt docker-cfg/ca-certificates.crt

test: $(BASE)/tmp/dynamo.containerid $(BASE)/tmp/redis.containerid jstest gotest pytest dockerstop

newvendor:
	govendor add +external
