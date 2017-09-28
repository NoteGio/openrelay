PACKAGE  = github.com/notegio/openrelay
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)
GOSTATIC = go build -a -installsuffix cgo -ldflags '-extldflags "-static"'


all: bin nodesetup truffleCompile

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

clean:
	rm -rf bin/ .gopath/

nodesetup:
	cd $(BASE)/js ; npm install

bin/delayrelay: $(BASE) cmd/delayrelay/main.go
	cd $(BASE) &&  CGO_ENABLED=0 $(GOSTATIC) -o bin/delayrelay cmd/delayrelay/main.go

bin/fundcheckrelay: $(BASE) cmd/fundcheckrelay/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/fundcheckrelay cmd/fundcheckrelay/main.go

bin/getbalance: $(BASE) cmd/getbalance/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/getbalance cmd/getbalance/main.go

bin/ingest: $(BASE) cmd/ingest/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/ingest cmd/ingest/main.go

bin/initialize: $(BASE) cmd/initialize/main.go
	cd $(BASE) && CGO_ENABLED=0 $(GOSTATIC) -o bin/initialize cmd/initialize/main.go

bin/simplerelay: $(BASE) cmd/simplerelay/main.go
	cd $(BASE) && CGO_ENABLED=0 $(GOSTATIC) -o bin/simplerelay cmd/simplerelay/main.go

bin/validateorder: $(BASE) cmd/validateorder/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/validateorder cmd/validateorder/main.go

bin: bin/delayrelay bin/fundcheckrelay bin/getbalance bin/ingest bin/initialize bin/simplerelay bin/validateorder

truffleCompile:
	cd $(BASE)/js ; node_modules/.bin/truffle compile

test:
	cd $(BASE)/accounts && go test
