# 1.目录结构

# 1.1 各个模块介绍
|     目录      |    子目录     |               说明               |                                   备注                                    |
|:-----------:|:----------:|:------------------------------:|:-----------------------------------------------------------------------:|
|     cmd     |    api     |            api 服务代码            |                       包含http server 和 RPC client                        |
|             |  comment   |         comment 微服务代码          |                                                                         |
|             |  favorite  |         favorite 微服务代码         |                                                                         |
|             |  relation  |         relation 微服务代码         |                                                                         |
|             |    user    |           user 微服务代码           |                                                                         |
|             |   video    |          video 微服务代码           |                                                                         |
|   config    |            |         微服务以及第三方包的配置文件         |                                                                         |
|     dal     |     db     |          操作 MySQL 代码           |                   包含 Gorm 初始化、Gorm 结构体以及 MySQL 的操作逻辑                    |
|             |   redis    |          操作 Redis 的代码          |       包含 go-redis 初始化、Redis 结构体、Redis 操作逻辑以及 Redis 与MySQL 数据同步的逻辑       |
|  internal   |  response  |            封装返回的结构体            |                                                                         |
|             |    tool    |          封装项目中常用的工具函数          |               包含 RSA、MD5、SHA256 等加密算法以及操作 ffmpeg 对上传视频截图                |
|    kitex    |            | 包含 proto 文件以及由 Kitex 生成的 go 代码 |                     Kitex 生成的 Go 代码在 Kitex_gen 子目录下                     |
|     pic     |            |             保存一些图片             |                                                                         |
|     pkg     |   errno    |             错误码的实现             |                                                                         |
|             |    etcd    |            服务注册与发现             |                                                                         |
|             |   gocron   |              定时任务              |                                                                         |
|             |    jwt     |           token生成与校验           |                                                                         |
|             | middelware |             常用的中间件             |                                                                         |
|             |   minio    |              对象存储              |                                                                         |
|             |  rabbitmq  |              消息队列              |                                                                         |
|             |   viper    |              配置读取              |                                                                         |
|             |    zap     |           日志打印与日志切割            |                                                                         |
|   scripts   |            |             存放启动脚本             | 包含使用 Docker 启动的 MySQL、Redis、etcd、rabbitmq、minio 等 |

树形结构
```bash
.
├── LICENSE
├── README.md
├── cmd
│   ├── api
│   │   ├── handler
│   │   │   ├── comment.go
│   │   │   ├── favorite.go
│   │   │   ├── relation.go
│   │   │   ├── user.go
│   │   │   └── video.go
│   │   ├── main.go
│   │   └── rpc
│   │       ├── comment.go
│   │       ├── favorite.go
│   │       ├── init.go
│   │       ├── relation.go
│   │       ├── user.go
│   │       └── video.go
│   ├── comment
│   │   ├── build.sh
│   │   ├── main.go
│   │   ├── script
│   │   │   └── bootstrap.sh
│   │   └── service
│   │       ├── handler.go
│   │       └── init.go
│   ├── favorite
│   │   ├── build.sh
│   │   ├── main.go
│   │   ├── script
│   │   │   └── bootstrap.sh
│   │   └── service
│   │       ├── handler.go
│   │       ├── init.go
│   │       └── timer.go
│   ├── relation
│   │   ├── build.sh
│   │   ├── main.go
│   │   ├── script
│   │   │   └── bootstrap.sh
│   │   └── service
│   │       ├── handler.go
│   │       ├── init.go
│   │       └── timer.go
│   ├── user
│   │   ├── build.sh
│   │   ├── main.go
│   │   ├── script
│   │   │   └── bootstrap.sh
│   │   └── service
│   │       ├── handler.go
│   │       └── init.go
│   └── video
│       ├── build.sh
│       ├── main.go
│       ├── script
│       │   └── bootstrap.sh
│       └── service
│           ├── handler.go
│           ├── init.go
│           └── upload.go
├── config
│   ├── api.yml
│   ├── comment.yml
│   ├── crypt.yml
│   ├── db.yml
│   ├── favorite.yml
│   ├── log.yml
│   ├── minio.yml
│   ├── rabbitmq.yml
│   ├── relation.yml
│   ├── user.yml
│   └── video.yml
├── dal
│   ├── db
│   │   ├── comment.go
│   │   ├── favorite.go
│   │   ├── feed.go
│   │   ├── init.go
│   │   ├── publish.go
│   │   ├── relation.go
│   │   └── user.go
│   └── redis
│       ├── common.go
│       ├── favorite.go
│       ├── init.go
│       └── relation.go
├── go.mod
├── go.sum
├── internal
│   ├── response
│   │   ├── base.go
│   │   ├── comment.go
│   │   ├── favorite.go
│   │   ├── relation.go
│   │   ├── user.go
│   │   └── video.go
│   └── tool
│       ├── crypt.go
│       └── snapshot.go
├── kitex
│   ├── comment.proto
│   ├── favorite.proto
│   ├── kitex.sh
│   ├── kitex_gen
│   │   ├── comment
│   │   │   ├── comment.pb.fast.go
│   │   │   ├── comment.pb.go
│   │   │   └── commentservice
│   │   │       ├── client.go
│   │   │       ├── commentservice.go
│   │   │       ├── invoker.go
│   │   │       └── server.go
│   │   ├── favorite
│   │   │   ├── favorite.pb.fast.go
│   │   │   ├── favorite.pb.go
│   │   │   └── favoriteservice
│   │   │       ├── client.go
│   │   │       ├── favoriteservice.go
│   │   │       ├── invoker.go
│   │   │       └── server.go
│   │   ├── relation
│   │   │   ├── relation.pb.fast.go
│   │   │   ├── relation.pb.go
│   │   │   └── relationservice
│   │   │       ├── client.go
│   │   │       ├── invoker.go
│   │   │       ├── relationservice.go
│   │   │       └── server.go
│   │   ├── user
│   │   │   ├── user.pb.fast.go
│   │   │   ├── user.pb.go
│   │   │   └── userservice
│   │   │       ├── client.go
│   │   │       ├── invoker.go
│   │   │       ├── server.go
│   │   │       └── userservice.go
│   │   └── video
│   │       ├── video.pb.fast.go
│   │       ├── video.pb.go
│   │       └── videoservice
│   │           ├── client.go
│   │           ├── invoker.go
│   │           ├── server.go
│   │           └── videoservice.go
│   ├── relation.proto
│   ├── user.proto
│   └── video.proto
├── pic
│   ├── 后端架构图.png
├── pkg
│   ├── errno
│   │   ├── code.go
│   │   └── errno.go
│   ├── etcd
│   │   ├── common.go
│   │   ├── discovery.go
│   │   └── registry.go
│   ├── gocron
│   │   ├── gocron.go
│   │   └── gocron_test.go
│   ├── jwt
│   │   ├── jwt.go
│   │   └── jwt_test.go
│   ├── middleware
│   │   ├── client.go
│   │   ├── common.go
│   │   ├── server.go
│   ├── minio
│   │   ├── init.go
│   │   ├── minio.go
│   ├── rabbitmq
│   │   ├── emit_logs.go
│   │   ├── init.go
│   │   ├── rabbitmq.go
│   │   └── receive_logs.go
│   ├── viper
│   │   └── viper.go
│   └── zap
│       └── zap.go
├── scripts
├── etcd
├── minio
├── mysql
├── rabbitmq
├── redis
└── docker-compose.yml


```


# 1.2项目启动
```bash
cd scripts
docker-compose up -d
```
然后启动cmd文件夹内个微服务文件夹的主函数以及api文件夹的主函数
(mysql需要手动创建db_tiktok数据库,该数据库名在config/db.tml中定义)
---

# 2.项目实现
## 2.1 技术选型与相关开发文档
本项目包含三大类接口：基础接口、互动接口、社交接口。采用微服务架构以及 Docker 部署的方式。总共需要 4G 存储空间，1 台服务器，项目中所需要的数据库以及中间件均由 Docker 下载并挂载运行。

以下是开发文档。

https://bytedance.feishu.cn/docx/BhEgdmoI3ozdBJxly71cd30vnRc

## 2.2 架构设计

由于抖音的各接口之间耦合度不高，故采用微服务的架构，分为 HTTP、RPC 和 DAL 三层：

- 其中 HTTP 层使用 Hertz 框架接收客户端发来的 HTTP 请求；
- RPC 层使用 Kitex 框架，并用 ETCD 做服务注册与发现；
- DAL 层为数据访问层，包含 Redis 和 MySQL 两部份。

当客户端发来 HTTP 请求时，HTTP 层会调用 RPC层 的 RPC Client，然后 RPC Client 去 ETCD 中心寻找已经注册的对应的微服务，交给
RPC Server 处理。RPC Server 会去调用 DAL 层的数据库，数据库处理完毕后把结果返回给 RPC Server，RPC Server 将其返回给 RPC
Client，最后 RPC Client 返回给 Hertz 并由 Hertz 返回 HTTP 响应结果。

![5a0c9126c565ce67bbb3b7fd9e1968fe.png](https://x.imgs.ovh/x/2023/09/03/64f3c7db3ecee.png)

# 3.测试结果

## 3.1 功能测试

|  功能项  | 功能需求      | 测试点                                          |        模块        |  结果  |
|:-----:|-----------|----------------------------------------------|:----------------:|:----:|
| 基础功能项 | 视频 Feed 流 | 支持所有用户刷抖音，视频按投稿时间倒序推出                        |      获取视频列表      | 正确运行 |
|       | 视频投稿      | 支持登录用户自己拍视频投稿                                |       发布视频       | 正确运行 |
|       | 个人主页      | 支持查看用户基本信息和投稿列表，注册用户流程简化                     |        注册        | 正确运行 |
|       |           |                                              |        登录        | 正确运行 |
|       |           |                                              |       个人信息       | 正确运行 |
| 方向功能项 | 喜欢列表      | 登录用户可以对视频点赞，在个人主页喜欢Tab 下能够查看点赞视频列表           |      获取喜欢列表      | 正确运行 |
|       |           |                                              |        点赞        |  正确运行|
|       |           |                                              |       取消赞        | 正确运行 |
|       | 用户评论      | 支持未登录用户查看视频下的评论列表，登录用户能够发表评论                 |      获取评论列表      | 正确运行 |
|       |           |                                              |       新增评论       | 正确运行 |
|       |           |                                              |       删除评论       | 正确运行 |
|       | 关系列表      | 登录用户可以关注其他用户，能够在个人主页查看本人的关注数和粉丝数，查看关注列表和粉丝列表 |        关注        |  正确运行|
|       |           |                                              |        取关        | 正确运行 |
|       |           |                                              | 获取关系列表（关注、粉丝、朋友） |  正确运行|
