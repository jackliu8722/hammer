# 巴别塔后端

## 需要
1. go 1.8+
1. dep 包管理器
1. mysql
1. redis
1. [eos-go](https://github.com/cookedsteak/eos-go)
1. eos节点(开放api插件）

## 目录结构
```console
├── hammer
│   ├── config      配置文件
│   ├── contract    合约
│   ├── controller  控制器
│   ├── logic       业务逻辑
│   ├── middleware  中间件
│   ├── model       数据库orm
│   ├── server      服务启动入口
│   ├── util        公共库与服务
│   ├── vendor      包
```
