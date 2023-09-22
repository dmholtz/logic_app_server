# logic_app_server

## Install

Install all dependencies using `go mod download`.

## Run

To run the server, first create a database. This needs to be done only on the first start-up:

```bash
go run cmd/db_setup/main.go
```

To start the server, run:

```bash
go run cmd/server/main.go
```
