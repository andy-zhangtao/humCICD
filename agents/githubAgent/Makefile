
.PHONY: build
name = github-agent

build:
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

run: build
	./$(name)

release: *.go
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o $(name)
	docker build -t vikings/$(name) .
	dpush -i vikings/$(name)
	rm $(name)