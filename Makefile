srcs = main.go
PROG = bin/unicorn
GOCMD = go
CC = build
CFLAGS = -x -ldflags "-s -w" -gcflags "-l"

all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	$(GOCMD) $(CC) -o $(PROG) $(CFLAGS)

alpine:
		if [ ! -d "./bin/" ]; then \
    	mkdir bin; \
    	fi
    	$(GOCMD) $(CC) -o $(PROG) -a $(CFLAGS)

install:
	cp bin/unicorn $GOPATH/bin

clean:
	rm -rf bin