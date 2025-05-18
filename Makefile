.PHONY: help

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

all: help

.PHONY:
help:
	@echo "go-dengueboy"
	@echo ""
	@echo "run -> Build the app and run"
	@echo ""

.PHONY:
run: 
	go run cmd/main.go