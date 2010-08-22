.PHONY: all clean install format

all: obsidian

include ${GOROOT}/src/Make.${GOARCH}
include Makefile.info

.SUFFIXES: .go .${O}

obsidian: main.${O}
	${LD} -o $@ main.${O}

.go.${O}:
	${GC} $*.go

.go.a:
	${GC} -o $*.${O} $*.go && gopack grc $*.a $*.${O}

format:
	gofmt -w ${GOFILES}

clean:
	rm -f obsidian ${GOPACKAGES} ${GOARCHIVES}

include Makefile.deps
