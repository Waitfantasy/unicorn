srcs = main.go
PROG = bin/unicorn
GOCMD = go
CC = build
CFLAGS = -x -ldflags "-s -w"
all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	$(GOCMD) $(CC) -o $(PROG) $(CFLAGS)

alpine: export CGO_ENABLED=0

alpine:
	if [ ! -d "./bin/" ]; then \
    		mkdir bin; \
    	fi
	$(GOCMD) $(CC) -o $(PROG) -a $(CFLAGS)
install:
	cp bin/unicorn $GOPATH/bin

clean:
	rm -rf bin
