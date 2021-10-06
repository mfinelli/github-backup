SOURCES := $(wildcard *.go)

all: ghb

clean:
	rm -f ghb

ghb: $(SOURCES)
	go build -o $@ $(SOURCES)

.PHONY: all clean
