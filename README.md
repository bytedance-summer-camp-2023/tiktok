# 目录结构
```bash
|-- cmd
|   |-- api
|   |   |-- handler
|   |   `-- rpc
|   `-- user
|       `-- service
|-- config
|-- dal
|   `-- db
|-- internal
|   |-- response
|   `-- tool
|-- kitex
|   `-- kitex_gen
|       `-- user
|           `-- userservice
|-- pkg
|   |-- viper
|   `-- zap
`-- scripts
```

---

# 项目启动
```bash
cd scripts
docker-compose up -d
```
然后启动各个微服务

---

# 项目实现
## 2.1 技术选型与相关开发文档
本项目包含三大类接口：基础接口、互动接口、社交接口。采用微服务架构以及 Docker 部署的方式。总共需要 16G 存储空间，1 台服务器，项目中所需要的数据库以及中间件均由 Docker 下载并挂载运行。

以下是开发文档。

https://bytedance.feishu.cn/docx/BhEgdmoI3ozdBJxly71cd30vnRc

## 2.2 架构设计


# 测试结果

## 3.1 功能测试

|  功能项  | 功能需求      | 测试点                                          |        模块        |  结果  |
|:-----:|-----------|----------------------------------------------|:----------------:|:----:|
| 基础功能项 | 视频 Feed 流 | 支持所有用户刷抖音，视频按投稿时间倒序推出                        |      获取视频列表      |  |
|       | 视频投稿      | 支持登录用户自己拍视频投稿                                |       发布视频       |  |
|       | 个人主页      | 支持查看用户基本信息和投稿列表，注册用户流程简化                     |        注册        | 正确运行 |
|       |           |                                              |        登录        | 正确运行 |
|       |           |                                              |       个人信息       | 正确运行 |
| 方向功能项 | 喜欢列表      | 登录用户可以对视频点赞，在个人主页喜欢Tab 下能够查看点赞视频列表           |      获取喜欢列表      |  |
|       |           |                                              |        点赞        |  |
|       |           |                                              |       取消赞        |  |
|       | 用户评论      | 支持未登录用户查看视频下的评论列表，登录用户能够发表评论                 |      获取评论列表      |  |
|       |           |                                              |       新增评论       |  |
|       |           |                                              |       删除评论       |  |
|       | 关系列表      | 登录用户可以关注其他用户，能够在个人主页查看本人的关注数和粉丝数，查看关注列表和粉丝列表 |        关注        |  |
|       |           |                                              |        取关        |  |
|       |           |                                              | 获取关系列表（关注、粉丝、朋友） |  |