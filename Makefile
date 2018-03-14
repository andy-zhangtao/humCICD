
.PHONY: build
name = hicd

build-trafficAgent:
	cd agents/trafficAgent; make

build-goAgent:
	cd agents/lanAgents/golang; make

build-buildAgent:
	cd agents/buildAgent; make

build-gitAgent:
	cd agents/gitAgent; make

build-client: build-goAgent build-buildAgent build-gitAgent build-trafficAgent
	echo "Build Agents"

build: build-client
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

all: build
	@echo "Build HICD"
    mv $(name) bin

run: build
	./$(name)

release: *.go *.md
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o $(name)
	mv $(name) bin/$(name)
	docker build -t vikings/$(name) .
	docker push vikings/$(name)
