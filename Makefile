setup:
	go install github.com/a-h/templ/cmd/templ@latest

build:
	templ generate
	go build -o ./bin/pastepass

test:
	go test ./...

run:
	./bin/pastepass
