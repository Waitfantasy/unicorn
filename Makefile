PROG =  ./bin/unicorn ./cmd
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')
CFLAGS = -ldflags "-s -w -X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\""
all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	go build $(CFLAGS) -o $(PROG) $(SRCS)

alpine: export CGO_ENABLED=0

alpine:
	if [ ! -d "./bin/" ]; then \
    		mkdir bin; \
	fi
	go build -o $(PROG) -a $(CFLAGS) $(SRCS)

install:
	cp ./cmd/unicorn/unicorn $GOPATH/bin

clean:
	@rm -rf bin
