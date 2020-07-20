.PHONY: all install

BINFILE = tttgameserver
MAINFILE = main.go

all: build

build:
	go build -o ${BINFILE} cmd/${MAINFILE}
