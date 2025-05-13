.PHONY: release import plan init install test_acc test docs clean mrproper

all: build test install

release:
	goreleaser release --clean

version = $(shell git describe --tags --abbrev=0)
os = $(shell uname -s|tr '[:upper:]' '[:lower:]')
arch = $(shell uname -m)

builds = \
	build/darwin_${arch}/terraform-provider-transip_${version} \
	build/linux_${arch}/terraform-provider-transip_${version}

build: ${builds}

# import test resources
import: init
	tofu import -config examples/ transip_domain.test $$TF_VAR_domain
	tofu import -config examples/ transip_vps.test $$TF_VAR_vps_name
	tofu state list | xargs -n1 tofu state show

comma=,
_targets = $(addprefix -target=,$(subst ${comma}, ,${targets}))

terraform = tofu -chdir=examples/

apply: init
	${terraform} apply ${_targets}

destroy: init
	${terraform} destroy ${_targets}

plan: init
	${terraform} plan -detailed-exitcode ${_targets}

init: build/terraform-provider-transip
	${terraform} init

dev_install: build/terraform-provider-transip
build/terraform-provider-transip:  build/${os}_${arch}/terraform-provider-transip_${version}
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
	@echo 'provider "aequitas/transip" {}' > tmp.tf
	mkdir -p docs/{resources,data-sources}/
	${terraform} providers schema -json | ./tools/docs.py
	@rm -f tmp.tf

clean:
	rm -rf terraform-provider-transip* build/ docs/
	rm -rf .terraform/ terraform.tfstate

mrproper: clean
	go clean -modcache
	rm -rf ./gopath/
