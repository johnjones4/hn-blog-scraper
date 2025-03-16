PROJECT=$(shell basename $(shell pwd))
TAG=ghcr.io/johnjones4/${PROJECT}
VERSION=$(shell date +%s)

all: info container

info:
	echo ${PROJECT} ${VERSION}

container:
	docker build --platform linux/x86_64 -t ${TAG} .
	docker push ${TAG}:latest
	docker image rm ${TAG}:latest

