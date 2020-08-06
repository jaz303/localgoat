SRC := $(shell find . -name '*.go')

lg: $(SRC)
	go build -o $@

clean:
	rm -f lg
