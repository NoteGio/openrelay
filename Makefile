PACKAGE  = github.com/notegio/openrelay
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)
GOSTATIC = go build -a -installsuffix cgo -ldflags '-extldflags "-static"'

env :
	export CGO_ENABLED=0
	export GOOS=linux

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

clean:
	rm -rf bin/ .gopath/

bin/delayrelay: $(BASE) env cmd/delayrelay/main.go
	cd $(BASE) &&  CGO_ENABLED=0 $(GOSTATIC) -o bin/delayrelay cmd/delayrelay/main.go

bin/fundcheckrelay: $(BASE) env cmd/fundcheckrelay/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/fundcheckrelay cmd/fundcheckrelay/main.go

bin/getbalance: $(BASE) env cmd/getbalance/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/getbalance cmd/getbalance/main.go

bin/ingest: $(BASE) env cmd/ingest/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/ingest cmd/ingest/main.go

bin/initialize: $(BASE) env cmd/initialize/main.go
	cd $(BASE) && CGO_ENABLED=0 $(GOSTATIC) -o bin/initialize cmd/initialize/main.go

bin/simplerelay: $(BASE) env cmd/simplerelay/main.go
	cd $(BASE) && CGO_ENABLED=0 $(GOSTATIC) -o bin/simplerelay cmd/simplerelay/main.go

bin/validateorder: $(BASE) env cmd/validateorder/main.go
	cd $(BASE) && $(GOSTATIC) -o bin/validateorder cmd/validateorder/main.go

bin: bin/delayrelay bin/fundcheckrelay bin/getbalance bin/ingest bin/initialize bin/simplerelay bin/validateorder
