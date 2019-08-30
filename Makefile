DEST=build
IMG=somefive/crd-discovery-server
TAG=0.1.0
syncrd:
	mkdir -p build/
	go build -o $(DEST)/syncrd github.com/Somefive/crd-discovery/cmd

serve:
	go build -o $(DEST)/serve github.com/Somefive/crd-discovery/server

install: syncrd
	sudo cp build/syncrd /usr/local/bin/kubectl-syncrd

clean:
	rm -rf build/

uninstall:
	sudo rm /usr/local/bin/kubectl-syncrd

docker-build:
	go mod vendor
	docker build . -t $(IMG):$(TAG) -f Dockerfile.server

docker-push: docker-build
	docker push $(IMG)