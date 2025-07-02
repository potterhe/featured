# featured

### init go mod

```
go mod init github.com/potterhe/featured
```

### cobra

[cobra-cli](https://github.com/spf13/cobra-cli/blob/main/README.md)

```
cobra-cli init --viper
```

### gRPC

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/helloworld/helloworld.proto
```

### buf

[buf](https://github.com/bufbuild/buf)

```
buf generate
```

`buf.yaml` Always run `buf dep update` after adding a dependency to your buf.yaml

```
buf dep update
```

### gRPC gateway

```
go get github.com/grpc-ecosystem/grpc-gateway/v2/runtime
```