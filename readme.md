### TCP转发服务

```go
编译
go build -o forward `

启动
forward -l 127.0.0.1:8080 -f 127.0.0.1:8081
```
