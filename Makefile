.PHONY: jet
jet:
	sqlite3 template.db < .\internal\repository\assets\init.sql
	jet -source=sqlite -dsn="template.db" -path=./internal/repository/.gen
	del template.db

.PHONY: debug
debug:
	go run cmd\main.go -debug

.PHONY: build
build:
	go build -ldflags "-s -w" -o mpwt.exe cmd\main.go
