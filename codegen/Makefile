all: run

run: build
	./codegen

build: codegen

codegen: *.go
	go build -o codegen

clean:
	rm -f codegen

.PHONY: all run build clean
