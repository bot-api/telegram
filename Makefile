.PHONY: all test test_v test_i lint vet fmt coverage checkfmt prepare errcheck race

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
PKGSDIRS=$(shell find -L . -type f -name "*.go" -not -path "./Godeps/*")

all: prepare

travis: checkfmt vet errcheck test_v race lint

prepare: fmt vet checkfmt errcheck test race lint

test_v:
	@echo "$(OK_COLOR)Test packages$(NO_COLOR)"
	@go test -cover -v ./...

test:
	@echo "$(OK_COLOR)Test packages$(NO_COLOR)"
	@go test -cover ./...

test_i:
ifdef API_BOT_TOKEN
	@echo "$(OK_COLOR)Run integration tests$(NO_COLOR)"
	@go test -cover -tags integration -run TestI_* ./...
endif

lint:
	@echo "$(OK_COLOR)Run lint$(NO_COLOR)"
	@test -z "$$(golint -min_confidence 0.3 ./... | grep -v Godeps/_workspace/src/ | tee /dev/stderr)"

vet:
	@echo "$(OK_COLOR)Run vet$(NO_COLOR)"
	@go vet ./...

errcheck:
	@echo "$(OK_COLOR)Run errchk$(NO_COLOR)"
	@errcheck

race:
	@echo "$(OK_COLOR)Test for races$(NO_COLOR)"
	@go test -race .

checkfmt:
	@echo "$(OK_COLOR)Check formats$(NO_COLOR)"
	@./scripts/checkfmt.sh .

fmt:
	@echo "$(OK_COLOR)Formatting$(NO_COLOR)"
	@echo $(PKGSDIRS) | xargs -I '{p}' -n1 goimports -w {p}

tools:
	@echo "$(OK_COLOR)Install tools$(NO_COLOR)"
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck
