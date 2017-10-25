.PHONY: all clean lint deps depspurge bin libs

SHELL=/bin/bash
LOCAL_DIRS = $(shell find * -maxdepth 0 -type d | grep -v '^vendor$$\|^bin$$')
LOCAL_PKGS = $(patsubst %, ./%/..., $(LOCAL_DIRS))
LOCAL_GOBIN = $(shell realpath -s $$PWD/bin)
LOCAL_NONGEN = $(shell find ${LOCAL_DIRS} -type f -iname '*.go' -a '!' -iname '*.capnp.go')

all: deps bin

clean:
	rm -f gocover.html vendor/.deps.stamp
	GOBIN=${LOCAL_GOBIN} go clean -i ${LOCAL_PKGS}
	go clean -i ./vendor/...

lint:
	@echo "======> goimports"
	out=$$(goimports -d -local github.com/netsec-ethz ${LOCAL_NONGEN}); if [ -n "$$out" ]; then echo "$$out"; exit 1; fi
	@echo "======> gofmt"
	out=$$(gofmt -d -s ${LOCAL_DIRS}); if [ -n "$$out" ]; then echo "$$out"; exit 1; fi
	@echo "======> go vet"
	go vet ${LOCAL_PKGS}

deps: vendor/.deps.stamp

vendor/.deps.stamp: vendor/vendor.json
	@echo "$$(date -Iseconds) Remove unused deps"; \
	    govendor list -no-status +unused | while read pkg; do \
	    grep -q '"path": "'$$pkg'"' vendor/vendor.json && continue; \
	    echo "$$pkg"; \
	    govendor remove "$$pkg"; \
	done
	@echo "$$(date -Iseconds) Syncing deps"; govendor sync -v
	@echo "$$(date -Iseconds) Installing deps"; go install ./vendor/...
	@if [ -n "$$(govendor list -no-status +outside | grep -v '^context$$')" ]; then \
	    echo "ERROR: external/missing packages:"; \
	    govendor list +outside; \
	    exit 1; \
	fi;
	touch $@

depspurge:
	rm -f vendor/.deps.stamp
	go clean -i ./vendor/...
	find vendor/* -maxdepth 0 -type d -exec rm -rf ./{} \;

bin: deps
	GOBIN=${LOCAL_GOBIN} govendor install -v +local,program

libs: deps
	govendor install -v +local,^program
