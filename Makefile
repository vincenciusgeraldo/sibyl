dep:
	go get -u github.com/kardianos/govendor
	govendor sync

bin:
	go build -o sibyl app/web/main.go

run: bin
	./sibyl