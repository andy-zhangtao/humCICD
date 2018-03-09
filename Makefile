
.PHONY: build
name = humcicd

build-goAgent:
	cd agents/lanAgents/golang; make

build-buildAgent:
	cd agents/buildAgent; make

build-gitAgent:
	cd agents/gitAgent; make

build-Agent:
	cd agents/agent; make

build-client: build-goAgent build-buildAgent build-gitAgent build-Agent
	echo "Build Agents"

build: build-client
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

all: build
	@echo "Build HICD"

run: build
	./$(name)

release: *.go *.md
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o $(name)
	docker build -t vikings/$(name) .
	docker push vikings/$(name)
