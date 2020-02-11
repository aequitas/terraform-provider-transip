version = $(shell git describe --tags --abbrev=0)

all: test release

releases = \
	terraform-provider-transip_${version}_darwin_amd64.tgz \
	terraform-provider-transip_${version}_linux_amd64.tgz

release: ${releases}

terraform-provider-transip_${version}_%_amd64.tgz: build/%_amd64/terraform-provider-transip_v${version}
	tar -zcf $@ -C ${<D} ${<F}

build/%_amd64/terraform-provider-transip_v${version}: $(wildcard *.go)
	mkdir -p ${@D}; GOOS=$* go build -o $@

test_integration: test
test_integration: TF_ACC=1

test:
	TF_ACC=${TF_ACC} go test -v

clean:
	rm -rf terraform-provider-transip* build/

mrproper: clean
	go clean -modcache
	rm -rf ./gopath/

.PHONY: release
