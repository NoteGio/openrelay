PACKAGE  = github.com/notegio/openrelay
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

bin/delayrelay: $(BASE) cmd/delayrelay/main.go
	cd $(BASE) && go build -o bin/delayrelay cmd/delayrelay/main.go

bin/fundcheckrelay: $(BASE) cmd/fundcheckrelay/main.go
	cd $(BASE) && go build -o bin/fundcheckrelay cmd/fundcheckrelay/main.go

bin/getbalance: $(BASE) cmd/getbalance/main.go
	cd $(BASE) && go build -o bin/getbalance cmd/getbalance/main.go

bin/ingest: $(BASE) cmd/ingest/main.go
	cd $(BASE) && go build -o bin/ingest cmd/ingest/main.go

bin/initialize: $(BASE) cmd/initialize/main.go
	cd $(BASE) && go build -o bin/initialize cmd/initialize/main.go

bin/simplerelay: $(BASE) cmd/simplerelay/main.go
	cd $(BASE) && go build -o bin/simplerelay cmd/simplerelay/main.go

bin/validateorder: $(BASE) cmd/validateorder/main.go
	cd $(BASE) && go build -o bin/validateorder cmd/validateorder/main.go

bin: bin/delayrelay bin/fundcheckrelay bin/getbalance bin/ingest bin/initialize bin/simplerelay bin/validateorder
