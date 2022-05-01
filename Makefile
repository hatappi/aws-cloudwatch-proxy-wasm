export GOROOT=$(shell go env GOROOT)

build-all: build-sender build-receiver

build-sender:
	tinygo build -o sender.wasm -scheduler=none -target=wasi ./cmd/sender/main.go

build-receiver:
	tinygo build -o receiver.wasm -scheduler=none -target=wasi ./cmd/receiver/main.go
