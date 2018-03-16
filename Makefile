
.PHONY: build
name = hicd

build-trafficAgent:
	cd agents/trafficAgent; make

release-trafficAgent:
	cd agents/trafficAgent; make release

build-goAgent:
	cd agents/lanAgents/golang; make

release-goAgent:
	cd agents/lanAgents/golang; make release

build-buildAgent:
	cd agents/buildAgent; make

release-buildAgent:
	cd agents/buildAgent; make release

build-gitAgent:
	cd agents/gitAgent; make

release-gitAgent:
	cd agents/gitAgent; make release

build-client: build-goAgent build-buildAgent build-gitAgent build-trafficAgent
	echo "Build Agents"

release-client: release-goAgent release-buildAgent release-gitAgent release-trafficAgent
	echo "Release Agents"

build: build-client
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

release-hicd: *.go *.md
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o $(name)
	mv $(name) bin/$(name)
	docker build -t vikings/$(name) .
	docker push vikings/$(name)

release: release-client

all: build release release-hicd
	@echo "Build HICD"
	mv $(name) bin


run: build
	./$(name)


