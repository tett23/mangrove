SHELL := /usr/local/bin/zsh

build-bindata: assets/bindata.go
	go-bindata -o assets/bindata.go -pkg assets config
