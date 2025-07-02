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