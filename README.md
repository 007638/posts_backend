## 帖子项目
项目简介：通过 Go + Gin + GORM + PostgreSQL 搭建后端接口、React 构建前端页面，实现用户注册登录、发帖、评论回复、帖子列表与个人中心等完整功能，适用于个人博客、学习论坛等场景。

## 技术栈
- Go
- Gin（Web 框架）
- GORM（ORM 数据库操作）
- PostgreSQL（数据库）

## 已实现接口
| 方法 | 路径 | 功能 |
| --- | --- | --- |
| POST | /api/register | 用户注册 |
| POST | /api/login | 用户登录 |
| GET | /api/posts | 帖子列表 |
| POST | /api/create | 发布帖子 |
| GET | /api/post/:id | 帖子详情（含评论） |
| POST | /api/comment | 发表评论/回复 |
| GET | /api/profile | 个人信息（帖子数/回帖数） |
| GET | /api/my/posts | 我的帖子 |
| GET | /api/my/comments | 我的回帖 |

## 安装步骤
1. 克隆项目到本地：
git clone https://github.com/007638/posts_backend.git

## 安装依赖并运行
cd posts_backend
go mod tidy
go run main.go

注意：运行前需要先启动PostgresSQL，并在main.go中配置好数据库连接
