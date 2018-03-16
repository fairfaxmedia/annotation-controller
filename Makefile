build: clean generate depend build-go

generate:
		git clone https://github.com/kubernetes/code-generator.git vendor/k8s.io/code-generator
		git clone https://github.com/kubernetes/apimachinery.git vendor/k8s.io/apimachinery
		./hack/update-codegen.sh

build-go:
		cd src/cmd/controller && go build

depend:
		dep ensure

clean: clean-vendor clean-go

clean-vendor:
		@rm -rf vendor
clean-go:
		go clean
image:
		docker build .
deploy:
		./scripts/create.sh
delete:
		./scripts/delete.sh