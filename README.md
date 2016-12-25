## grpc 扩展

进入 *.proto 所在的目录, 执行下面的命名
```
protoc --go_out=plugins=grpc+grpcx:. *.proto
```