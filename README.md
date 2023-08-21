Run docker-compose.yml and load mysql.ddl to your mysql.

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


# kitex生成文件
参考https://www.cloudwego.io/zh/docs/kitex/getting-started/

执行命令



# 项目实现
## 2.1 技术选型与相关开发文档
本项目包含三大类接口：基础接口、互动接口、社交接口。采用微服务架构以及 Docker 部署的方式。总共需要 16G 存储空间，1 台服务器，项目中所需要的数据库以及中间件均由 Docker 下载并挂载运行。

以下是开发文档。

https://bytedance.feishu.cn/docx/BhEgdmoI3ozdBJxly71cd30vnRc

## 2.2 架构设计
