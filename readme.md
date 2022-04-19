### TCP转发服务

```go
编译
go build -o forward `

启动
forward start --config=./config.yml

后台启动
nohup ./forward start --config=./config.yml  > ./runoob.log 2>&1 &
```
