.PHONY: release release_upload import plan init install test_acc test docs clean mrproper

version = $(shell git describe --tags --abbrev=0)
os = $(shell uname -s|tr '[:upper:]' '[:lower:]')
arch = amd64

release_ext ?= zip

SIGNING_ID ?= ED49BB55

all: test install

releases = \
	terraform-provider-transip_${version}_darwin_${arch}.${release_ext} \
	terraform-provider-transip_${version}_linux_${arch}.${release_ext}

signatures = \
	terraform-provider-transip_${version}_SHA256SUMS \
	terraform-provider-transip_${version}_SHA256SUMS.sig

release_artifacts = ${releases} ${signatures}

upload_release: ${release_artifacts} | hub
	hub release create $(addprefix -a ,$^) ${version}

release: ${release_artifacts}

terraform-provider-transip_${version}_SHA256SUMS: ${releases}
	shasum -a 256 $^ > $@

terraform-provider-transip_${version}_SHA256SUMS.sig: %.sig: %
	gpg --detach-sign -u ${SIGNING_ID} $<

builds = \
	build/darwin_${arch}/terraform-provider-transip_${version} \
	build/linux_${arch}/terraform-provider-transip_${version}

build: ${builds}

# import test resources
import: init
	terraform import -config examples/ transip_domain.test $$TF_VAR_domain
	terraform import -config examples/ transip_vps.test $$TF_VAR_vps_name
	terraform state list | xargs -n1 terraform state show

comma=,
_targets = $(addprefix -target=,$(subst ${comma}, ,${targets}))

apply: init
	terraform apply ${_targets} examples/

destroy: init
	terraform destroy ${_targets} examples/

plan: init
	terraform plan -detailed-exitcode ${_targets} examples/

init: .terraform/plugins/darwin_amd64/lock.json

.terraform/plugins/darwin_amd64/lock.json: terraform.d/plugins/${os}_${arch}/terraform-provider-transip_${version} | terraform
	terraform init examples/

install: terraform.d/plugins/${os}_${arch}/terraform-provider-transip_${version}
terraform.d/plugins/${os}_${arch}/terraform-provider-transip_${version}:  build/${os}_${arch}/terraform-provider-transip_${version}
	mkdir -p ${@D}
	cp $< $@

terraform-provider-transip_${version}_%_${arch}.zip: build/%_${arch}/terraform-provider-transip_${version}
	zip $@ $<

terraform-provider-transip_${version}_%_${arch}.tgz: build/%_${arch}/terraform-provider-transip_${version}
	tar -zcf $@ -C ${<D} ${<F}

build/%_${arch}/terraform-provider-transip_${version}: $(wildcard *.go) go.mod
	mkdir -p ${@D}; GOOS=$* go build -o $@

test_acc: test
test_acc: TF_ACC=1

test:
	TF_ACC=${TF_ACC} go test -v

docs: | init
	@echo 'provider "transip" {}' > tmp.tf
	mkdir -p docs/{resources,data-sources}/
	terraform providers schema -json | ./tools/docs.py
	@rm -f tmp.tf

clean:
	rm -rf terraform-provider-transip* build/ docs/
	rm -rf .terraform/ terraform.tfstate

mrproper: clean
	go clean -modcache
	rm -rf ./gopath/

hub terraform: %: /usr/local/bin/%
/usr/local/bin/%:
	brew install $*