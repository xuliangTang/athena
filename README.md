## 关于 Athena

基于 `Gin` 的轻量级套件脚手架，提供了开箱即用，强大丰富的基础组件库，有类似 `Spring Cloud` 框架灵活的注解、强大的全局依赖注入容器、标准的 PSR 规范实现等等。

- [x] 模块化、松耦合设计
- [x] 组件丰富、开箱即用
- [x] 稳健的工程设计规范
- [x] 容器管理，依赖注入（DI）
- [x] 强大的注解功能
- [x] 秒级定时任务和异步协程任务
- [x] 便携的CLI开发工具、自动代码生成
- [x] 自动化的测试用例生成

## 架构

![](https://s3.bmp.ovh/imgs/2022/11/02/389da2868270167e.png)

## 目录设计

业务项目基本目录结构如下示例

```
├── src/    ----- 应用代码目录
│   ├── beans/
│   │   └── Beans.go
│   ├── config/ 
│   ├── db/
│   ├── http/
│   │   ├── classes/
│   │   └── middlewares/
│   ├── interfaces/ 
│   ├── models/
│   │   ├── daos/
│   │   └── entity/
│   ├── request/
│   ├── service/ 
│   ├── tools/ 
│   ├── tasks/
│   ├── rpc/
│   │   └── services/
│   │   └── middlewares/
│   ├── webSocket/
│   │   ├── chats/
│   │   ├── middlewares/
│   ├── tcp/
│   │   └── classes/
│   └── main.go
├── storage/               ----- 临时文件目录（日志、上传文件、文件缓存等）
├── test/                  ----- 单元测试目录
├── application.yml
├── go.mod
```

## 安装

```
go get github.com/XNXKTech/athena
```

## 快速开始

只需简单的几行代码，就完成了一个 `http` 服务的构建

```go
package main

func main() {
	athena.Ignite(config.BaseConf).
		Load(NewConfigModule(),
			wechat.NewMiniProgramModule()).
		Beans(beans.Import()...).
		Attach(auth.NewAuthMiddleware()).
		Mount("v1", classes.NewEpisodeClass(),
			classes.NewUserClass()).
		CronTask("0/30 * * * * *", tasks.SyncEpisodes()).
		Launch()
}
```
