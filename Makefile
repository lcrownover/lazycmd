build:
	go build -o ./bin/lazycmd_server cmd/lazycmd/server/main.go

run: build
	./bin/lazycmd_server

clean:
	go clean
	rm -f ./bin/lazycmd_server
