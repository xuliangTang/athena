## 关于 Athena

基于 `Gin` 的轻量级套件脚手架，提供了开箱即用，强大丰富的基础组件库，有类似 `Spring Cloud` 框架灵活的注解、强大的全局依赖注入容器、标准的 PSR 规范实现等等。

- [x] 模块化、松耦合设计
- [x] 组件丰富、开箱即用
- [x] 稳健的工程设计规范
- [x] 容器管理，依赖注入（DI）
- [x] 强大的注解功能
- [x] 秒级定时任务和异步协程任务
- [x] 支持限流和熔断
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
├── athena-cli             ----- 命令行工具
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
	athena.Ignite().
		Load(NewConfigModule(),
			wechat.NewMiniProgramModule()).
		Beans(beans.Import()...).
		Attach(auth.NewAuthMiddleware()).
		Mount("v1", nil, classes.NewEpisodeClass(),
			classes.NewUserClass()).
		CronTask("0/30 * * * * *", tasks.SyncEpisodes()).
		Launch()
}
```

使用命令行工具生成一个 `Controller`
```bash
$ ./athena-cli new controller user
Controller [src/classes/UserClass.go] created successfully.
```

### 依赖注入
向`BeanFactory`中注册后，在需要的地方打上`inject`注解即可自动完成注入
```
type UserClass struct {
	UserService *service.UserService `inject:"-"`
}

type UserService struct {
	UserDAO      *daos.UserDAO   `inject:"-"`
	Db           *db.GormAdapter `inject:"-"`
	IpGeoService *IpGeoService   `inject:"-"`
}
```

### 限流
增加配置项：
```
rateLimitRules:
  - /v1/test:
      interval: 60  # 多长时间添加一次令牌
      capacity: 2   # 令牌桶的容量
      quantum: 2    # 到达定时器指定的时间，往桶里面加多少令牌
  - /v1/ping:
      interval: 1
      capacity: 1
      quantum: 2
```
加入中间件启用限流：
```
athena.Ignite().Attach(athena.NewRateLimit())
```

### 熔断
增加配置项：
```
fuseRules:
  - test1:                              # name
      timeout: 1000                     # command超时时间
      maxConcurrentRequests: 1000       # command最大并发量
      sleepWindow: 6000                 # 熔断时长
      requestVolumeThreshold: 1000      # 请求数量临界值
      errorPercentThreshold: 50         # 失败率阈值
```
启用配置：
```
athena.Ignite().Load(athena.NewFuse())
```
使用示例：
```
hystrix.Do("test1", func() error {
    resp, err := http.Get("https://www.google.com/")
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("请求失败:%v", err)
        return errors.New(fmt.Sprintf("error resp"))
    }
    return nil
    
}, func(err error) error {
    if err != nil {
        fmt.Printf("circuitBreaker and err is %s\n", err.Error())
        msg = err.Error()
    }
    return nil
})
```