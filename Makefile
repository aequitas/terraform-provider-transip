version = $(shell git describe --tags)

releases = \
	terraform-provider-transip_${version}_darwin_amd64.tgz \
	terraform-provider-transip_${version}_linux_amd64.tgz

release: ${releases}

terraform-provider-transip_${version}_%_amd64.tgz: build/%_amd64/terraform-provider-transip_v${version}
	tar -zcf $@ -C ${<D} ${<F}

build/%_amd64/terraform-provider-transip_v${version}: $(wildcard *.go)
	mkdir -p ${@D}; go build -o $@

clean:
	rm -rf terraform-provider-transip*.tgz build/

.PHONY: release
