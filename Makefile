version = $(shell git describe --tags --abbrev=0)
os = $(shell uname -s|tr '[:upper:]' '[:lower:]')
arch = amd64

release_ext ?= tgz

all: test install release

releases = \
	terraform-provider-transip_${version}_darwin_${arch}.${release_ext} \
	terraform-provider-transip_${version}_linux_${arch}.${release_ext}

release: ${releases}

builds = \
	build/darwin_${arch}/terraform-provider-transip_v${version} \
	build/linux_${arch}/terraform-provider-transip_v${version}

build: ${builds}

# import test resources
import: init
	terraform import -config examples/ transip_domain.test $$TF_VAR_domain
	terraform import -config examples/ transip_vps.test $$TF_VAR_vps_name
	terraform state list | xargs -n1 terraform state show

apply: init
	terraform apply -parallelism=1 examples/

plan: init
	terraform plan examples/

init: .terraform/plugins/darwin_amd64/lock.json

.terraform/plugins/darwin_amd64/lock.json: terraform.d/plugins/${os}_${arch}/terraform-provider-transip_v${version}
	terraform init examples/

install: terraform.d/plugins/${os}_${arch}/terraform-provider-transip_v${version}
terraform.d/plugins/${os}_${arch}/terraform-provider-transip_v${version}:  build/${os}_${arch}/terraform-provider-transip_v${version}
	mkdir -p ${@D}
	cp $< $@

terraform-provider-transip_${version}_%_${arch}.zip: build/%_${arch}/terraform-provider-transip_v${version}
	zip $@ $<

terraform-provider-transip_${version}_%_${arch}.tgz: build/%_${arch}/terraform-provider-transip_v${version}
	tar -zcf $@ -C ${<D} ${<F}

build/%_${arch}/terraform-provider-transip_v${version}: $(wildcard *.go)
	mkdir -p ${@D}; GOOS=$* go build -o $@

test_integration: test
test_integration: TF_ACC=1

test:
	TF_ACC=${TF_ACC} go test -v

clean:
	rm -rf terraform-provider-transip* build/
	rm -rf .terraform/ terraform.tfstate

mrproper: clean
	go clean -modcache
	rm -rf ./gopath/

.PHONY: release
