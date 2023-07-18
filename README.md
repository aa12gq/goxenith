# go-xenith

项目名称GoXenith(Xenith) 传达出顶点、巅峰的意义, 意将其打造为高性能、功能齐全的 API 框架, 以论坛社区为主题，基于gin, cobra, viper, zap, ent, protobuf, redis, mysql, email, jwt. 

# 项目wiki
https://github.com/AA12GQ/goxenith/wiki/

# 构建
在项目根目录下运行
```
    make all
```
# 部署

```
    GOOS=linux GOARCH=amd64 make build-direct
```

# 运行

```
    nohup ./goxenith &
```