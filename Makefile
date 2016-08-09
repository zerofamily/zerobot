.PHONY: build

build: 
	go build -o bin/zerobot cmd/main.go 
	

run: build
	./bin/zerobot -debug=true
