#compile go application

all: build

build: 
	go build -mod=vendor -ldflags="-w -s" -o main


.PHONY: all build
