dep:
	go get -u github.com/kardianos/govendor
	govendor sync

build-web:
	go build -o sibyl app/web/main.go

build-cron:
	go build -o sibyl-cron app/cron/main.go

run: build-web
	./sibyl

run-cron: build-cron
	./sibyl-cron