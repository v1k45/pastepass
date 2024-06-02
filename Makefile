setup:
	go install github.com/a-h/templ/cmd/templ@latest

build:
	templ generate
	go build -o ./bin/pastepass

run:
	./bin/pastepass
