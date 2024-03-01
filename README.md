Demo grpc service for trying out Mangle deductive database.

Mangle is an extension of Datalog.
The repo is here: https://github.com/google/mangle

## Building

```go
go get ./...
go build ./...
```

## Run the Server

This starts the gprc server, with a tiny database.

```
go run ./server --source=example/demo.mg
```

## Run the Client

This queries the gprc server, with a recursive query.

```go
go run ./client
go run ./client --query="reachable(X, /d)"
```

## Regenerate the proto files

```shell
cd proto
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative mangle.proto
```

