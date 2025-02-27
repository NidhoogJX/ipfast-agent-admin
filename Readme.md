# IPFAST-V1.0.0

## 目录
- [项目说明](#项目说明)
- [安装依赖](#安装依赖)
- [安装部署](#安装部署)
- [主要功能](#主要功能)

## 项目说明
IPFAST后台服务
提供 用户和相关后台接口
#### 目录结构
```
.
├── Dockerfile docker容器镜像打包文件 
├── Readme.md  项目文档说明
├── build.sh  镜像运行命令
├── cmd  项目入口
│   ├── api 后台服务
│   │   ├── config.yaml 配置文件 
│   │   ├── ipfast_group.yml Docker compose配置文件
│   │   ├── locales 其他资源文件
│   │   │   ├── admin 后台接口国际化
│   │   │   │   ├── en.json
│   │   │   │   ├── zh_CN.json
│   │   │   │   └── zh_TW.json
│   │   │   ├── en.json
│   │   │   ├── ip2location.csv IP库文件
│   │   │   ├── zh_CN.json
│   │   │   └── zh_TW.json
│   │   ├── log 日志
│   │   ├── main.go 服务入口
│   │   ├── main_test.go 测试服务入口
│   │   ├── script 数据库脚本
│   │   │   └── flowmanger.sql
│   │   └── static 静态资源文件 已废弃⚠️
│   │       └── country_icon
│   │           ├── CN.png
│   │           ├── DE.png
│   │           ├── ES.png
│   │           ├── FR.png
│   │           ├── GB.png
│   │           ├── IT.png
│   │           ├── RU.png
│   │           └── US.png
│   └── syncData
│       ├── config.yaml
│       ├── main.go
│       └── main_test.go
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── announcements
│   │   │   └── api.go
│   │   ├── commodites
│   │   │   └── api.go
│   │   ├── display
│   │   │   ├── api.go
│   │   │   └── country.go
│   │   ├── init.go
│   │   ├── pay
│   │   │   └── api.go
│   │   ├── proxy
│   │   │   └── api.go
│   │   └── user
│   │       └── api.go
│   ├── config
│   │   ├── api
│   │   │   └── init.go
│   │   ├── config.go
│   │   ├── countrycode
│   │   │   └── code.go
│   │   ├── i18n
│   │   │   └── i18n.go
│   │   ├── iplocation
│   │   │   └── ip.go
│   │   └── syncData
│   │       └── init.go
│   ├── db
│   │   ├── core
│   │   │   ├── gorm
│   │   │   │   └── gorm.go
│   │   │   ├── kafka
│   │   │   │   └── flow.go
│   │   │   └── redisHandler
│   │   │       └── redis.go
│   │   ├── library
│   │   │   ├── account.go
│   │   │   ├── user.go
│   │   │   └── verificationCodes.go
│   │   └── models
│   │       ├── account.go
│   │       ├── announcement.go
│   │       ├── assets.go
│   │       ├── commodites.go
│   │       ├── dataCenterIps.go
│   │       ├── durationType.go
│   │       ├── flowRecord.go
│   │       ├── ipRecord.go
│   │       ├── ipWhiteList.go
│   │       ├── ipipgoAccount.go
│   │       ├── ipipgoAccountStaticIp.go
│   │       ├── model.go
│   │       ├── payPlatform.go
│   │       ├── proxyServer.go
│   │       ├── regionCity.go
│   │       ├── regionCountry.go
│   │       ├── regionProvince.go
│   │       ├── staticIps.go
│   │       ├── subUser.go
│   │       ├── trafficConutryCommodites.go
│   │       ├── trafficCountry.go
│   │       ├── trafficRegion.go
│   │       ├── transactionOrders.go
│   │       ├── user.go
│   │       ├── userFlow.go
│   │       ├── verificationCode.go
│   │       └── xfj.go
│   ├── handler
│   │   ├── aesHandler
│   │   │   └── base.go
│   │   ├── emailHandler
│   │   │   └── base.go
│   │   ├── ginHandler
│   │   │   ├── aesAndGzipMiddleware.go
│   │   │   ├── base.go
│   │   │   ├── corsMiddleware.go
│   │   │   ├── jwtMiddleware.go
│   │   │   ├── pullMiddleware.go
│   │   │   ├── recaptchaMiddleware.go
│   │   │   └── translateMiddleware.go
│   │   ├── ipipgo
│   │   │   └── base.go
│   │   ├── network
│   │   │   ├── request
│   │   │   │   └── base.go
│   │   │   └── server
│   │   │       ├── base.go
│   │   │       └── response.go
│   │   └── stripeHandler
│   │       └── base.go
│   ├── scheduler
│   │   └── task.go
│   └── services
│       ├── announcement_services.go
│       ├── assets_services.go
│       ├── commodites_services.go
│       ├── duration_types_services.go
│       ├── ipipgo_services.go
│       ├── proxy_services.go
│       ├── static_data_serveices.go
│       ├── sync_traffic_services.go
│       ├── transaction_services.go
│       └── user_services.go
├── pkg
│   ├── aliyun
│   │   ├── captcha_v2
│   │   │   └── base.go
│   │   └── sms
│   │       └── base.go
│   ├── api
│   │   ├── accountFlow
│   │   │   └── accountFlow.pb.go
│   │   └── accountFlow.proto
│   ├── bloomFilter
│   │   └── base.go
│   └── util
│       ├── aliyun
│       ├── cronscheduler
│       │   └── base.go
│       ├── dingding
│       │   └── base.go
│       ├── geiTui
│       ├── ip
│       │   ├── ipdata
│       │   │   ├── city.csv
│       │   │   ├── country.csv
│       │   │   ├── ipv4.csv
│       │   │   ├── location.csv
│       │   │   └── province.csv
│       │   └── main.go
│       ├── log
│       │   └── log.go
│       ├── logBase
│       │   └── log.go
│       ├── messageKafka
│       │   └── base.go
│       └── spinnerHandler
│           └── base.go
└── test
    ├── check_pullorder_test.go
    ├── default.log
    ├── email_test.go
    ├── flowRecord_test.go
    ├── gomail_test.go
    ├── i18n_test.go
    ├── ipipgo_test.go
    ├── jwt_test.go
    ├── locales
    │   ├── admin
    │   │   ├── en.json
    │   │   ├── zh_CN.json
    │   │   └── zh_TW.json
    │   ├── en.json
    │   ├── zh_CN.json
    │   └── zh_TW.json
    ├── time_test.go
    └── token_test.go
```
## 安装依赖
```
go mod download

protobuf 格式生成
https://github.com/protocolbuffers/protobuf/releases 安装对应版本 添加到环境变量使用
vi ~/.zshrc or ~/.bash_profile
source ~/.bash_profile
source ~/.zshrc
下载项目 google.golang.org/protobuf/cmd/protoc-gen-go
进入 cmd/protoc-gen-go 目录下运行 go build
生成可执行二进制文件 丢到GOPATH/bin目录下
将 GOPATH/bin 添加到你的 PATH
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
protoc-gen-go --version 测试是否安装完成

protoc --version 测试是否正确安装
protoc --go_out=. device.proto #生成 Protocol协议对应的go消息结构文件
```
## 安装部署
```
两种部署方式
1.手动部署服务
2.Docker部署服务
```
#### 手动部署服务
```
1.准备阶段 打包 程序
cmd/dfj_load_balancer/ 目录下执行 交叉编译命令 具体根据系统
GOOS=linux GOARCH=amd64 go build -o ipfast
GOOS=darwin GOARCH=amd64 go build -o ipfast
将以下文件上传到服务器 同一个目录下 
cmd api里的资源文件和配置文件
2.后台运行服务 ToDo:后续添加为系统服务 添加开机自启动功能
chmod +x ipfast 给予执行权限
nohup ./ipfast > dfj.out &



```
#### Docker部署服务
```

1.上传部署文件
(1)安装Docker环境 (国内服务器需要修改镜像仓库源)
使用脚本 bash <(curl -sSL https://gitee.com/SuperManito/LinuxMirrors/raw/main/DockerInstallation.sh)
(2)执行命令启动所有服务:
docker compose -f ipfast_group.yml up -d
(3)停止并删除所有服务:
docker compose -f ipfast_group down --volumes # --rmi all 包括镜像
docker compose stop/start web 单独启动停止某个服务


```
#### 使用 Kafka 消息队列实现读写分离

1. **Kafka 简介**：Kafka 是由 Apache 软件基金会开发的一个开源流处理平台，由 Scala 和 Java 编写。该项目的目标是为处理实时数据提供一个统一、高吞吐、低延迟的平台。其持久化层本质上是一个“按照分布式事务日志架构的大规模发布/订阅消息队列”，这使它作为企业级基础设施来处理流式数据非常有价值。

2. **Kafka 的性能**：Kafka 的写入性能非常高，因为它是为处理大量实时数据流设计的。Kafka 使用了一些优化技术，例如顺序写入和零拷贝，这使得它能够以非常高的吞吐量写入数据。此外，Kafka 的数据是分布式存储的，你可以通过增加更多的分区和副本来提高写入性能。

3. **使用 Protobuf 协议**：消息队列的读写使用 Protobuf 协议。Protocol Buffers（Protobuf）是 Google 开发的一种数据序列化协议（类似于 XML、JSON、YAML 等），它能够将结构化数据序列化，可用于数据存储、通信协议等方面。Protobuf 相比 JSON，XML 格式的数据，数据更小（3 到 10 倍的压缩率）、速度更快（20 到 100 倍的速度），并且 Protobuf 提供了丰富的数据结构，并且可以生成各种语言的数据访问代码，包括 Go。


## 主要功能  [✅] [❌]
1.ToDo:待续


## pprof分析Go程序性能
```
go tool pprof http://localhost:6060/debug/pprof/heap  获取内存性能分析报告
go tool pprof http://localhost:6060/debug/pprof/profile 获取CPU性能分析报告
flat：函数直接分配的内存。  
flat%：函数直接分配的内存占总内存的百分比。  
sum%：到目前为止直接分配的内存占总内存的百分比。  
cum：函数及其所有子函数分配的内存。  
cum%：函数及其所有子函数分配的内存占总内存的百分比。

Commands:
    callgrind        输出callgrind格式的图表，这种格式可以被一些工具如KCachegrind或者QCacheGrind读取。
    comments         输出所有的分析注释。
    disasm           输出带有样本注释的汇编列表。
    dot              输出DOT格式的图表，这种格式可以被Graphviz等工具读取。
    eog              使用相应的工具可视化图表。eog
    evince           使用相应的工具可视化图表。 evince
    gif              输出相应格式的图表 GIF 
    gv               使用相应的工具可视化图表。 gv
    kcachegrind      在KCachegrind中可视化报告
    list             输出与正则表达式匹配的函数的带注释的源代码。
    pdf              输出相应格式的图表 PDF 
    peek             输出与正则表达式匹配的函数的调用者/被调用者。
    png              输出相应格式的图表 PNG 
    proto            输出压缩的protobuf格式的分析数据。
    ps               输出相应格式的图表 PS 
    raw              输出原始分析数据的文本表示。
    svg              输出相应格式的图表 SVG 
    tags             输出分析中的所有标签。
    text             以文本形式输出顶级条目。
    top              以文本形式输出顶级条目。
    topproto         以压缩的protobuf格式输出顶级条目。
    traces           以文本形式输出所有分析样本。
    tree             输出调用图的文本表示。
    web              在网页浏览器中可视化图表或显示带注释的源代码。
    weblist          在网页浏览器中可视化图表或显示带注释的源代码。
    o/options        列出所有选项及其当前值。
    q/quit/exit/^D   退出pprof
```