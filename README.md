# tiktok

### 项目介绍

抖声（抖音极速版）后端接口实现

### 功能说明

具体功能说明请参考飞书接口文档
https://bytedance.feishu.cn/docx/doxcnbgkMy2J0Y3E6ihqrvtHXPg

### 主要技术栈

- Gin
- Gorm
- MySQL
- 腾讯云对象存储COS
- JWT

### 项目环境

- Go 1.18
- MySQL 5.7
- Gorm 2.x

### 文件目录

~~~
├── controller （向客户端提供的api接口层)
│   ├── comment.go
│   ├── common.go (用于存放接口层共用的对象结构体)
│   ├── demo_data.go
│   ├── favorite.go
│   ├── feed.go
│   ├── publish.go
│   ├── relation.go
│   ├── stroage.go
│   └── user.go
├── docs （项目工作文档）
│   ├── Douyin-Database-v6 客户端.sql
│   ├── Douyin-Database-v6 服务器端.sql
│   ├── Gitee操作说明.md
│   ├── TencentCloudCOS.md
│   └── 更新说明.pdf
├── go.mod
├── go.sum
├── main.go
├── model （模型层和数据库操作层）
│   ├── dal.go
│   └── favorite.go
│   └── follow.go
│   └── mysql.go
│   └── user.go
│   └── video.go
├── pkg
│   ├── constants.go
│   └── errno.go
├── README.md
├── router.go
├── service （业务逻辑层）
│   ├── favorite.go
│   ├── publish.go
│   ├── relation.go
│   └── service.go
└── utils （工具类）
    ├── encrypt.go
    ├── jwt.go
    └── str2int.go