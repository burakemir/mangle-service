Demo grpc service for trying out Mangle deductive database.

## Building

```go
go get ./...
go build ./...
```

## Run the Server

```
go run ./server --source=example/demo.mg
```

## Run the Client

```go
go run ./client
go run ./client "reachable(X, /d)"
```

## Regenerate the proto files

```shell
cd proto
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative mangle.proto
```

