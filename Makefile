ARCH?=amd64
REPO?=#your repository here 
VERSION?=0.1

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -o ./bin/nats-subscriber main.go

container:
	docker build -t $(REPO)finops-nats-subscriber:$(VERSION) .
	docker push $(REPO)finops-nats-subscriber:$(VERSION)
