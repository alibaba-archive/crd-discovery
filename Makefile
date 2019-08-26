DEST=build
IMG=somefive/kube-sync-server
TAG=0.1.0
sync:
	go build -o $(DEST)/sync github.com/Somefive/crd-discovery/cmd

serve:
	go build -o $(DEST)/serve github.com/Somefive/crd-discovery/server

docker-build:
	go mod vendor
	docker build . -t $(IMG):$(TAG) -f Dockerfile.server

docker-push: docker-build
	docker push $(IMG)