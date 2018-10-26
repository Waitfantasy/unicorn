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

install:
	$(GOCMD) install

clean:
	rm -rf bin