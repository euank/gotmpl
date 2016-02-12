GO_FILES = $(shell find . -type f -name '*.go')

all: ./bin/gotmpl

./bin/gotmpl: $(GO_FILES)
	go build -o ./bin/gotmpl ./cmd/gotmpl

test: ./bin/gotmpl
	go test -v ./...

clean:
	rm -f ./bin/gotmpl
