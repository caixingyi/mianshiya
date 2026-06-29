# 面试刷题平台

Go + Next.js 全栈面试刷题平台。

## 技术栈

| 层级 | 技术 |
|---|---|
| 后端 | Go + Gin + GORM + Token Auth (Redis) |
| 前端 | Next.js 14 + Ant Design + ProComponents |
| 数据库 | MySQL + Redis + Elasticsearch |
| AI | 火山引擎 DeepSeek API |
| 部署 | Docker Compose + Nginx |

## 功能模块

- **用户系统**：注册/登录/个人信息编辑/头像上传/Redis Bitmap 签到
- **题库管理**：题库 CRUD、题目 CRUD、批量操作、题库-题目关联
- **AI 生成题目**：调用大模型自动生成面试题和题解（带重试机制）
- **AI 模拟面试**：多轮对话面试官，支持 start/chat/end 事件流
- **讨论区**：帖子发布/列表/详情、Markdown 编辑器、点赞/收藏
- **ES 全文搜索**：IK 中文分词、MySQL 降级、双写同步、gobreaker 熔断保护
- **限流防护**：Redis + Lua 固定窗口限流中间件
- **三级缓存**：进程内存(L1) + Redis(L2) + MySQL(L3)，热点自动发现与提升
- **配置管理**：Viper + YAML 配置，支持多环境

## 项目结构

```
mianshiya-go-backend/
  cmd/server/main.go       # API 服务入口
  internal/
    config/                # Viper 配置管理
    db/                    # MySQL/Redis 初始化
    auth/                  # Token 认证中间件
    user/                  # 用户模块
    question/              # 题目模块（含 AI 生成）
    questionbank/          # 题库模块（含三级缓存）
    post/                  # 帖子模块
    mockinterview/         # AI 模拟面试
    ai/                    # 火山引擎 AI 客户端
    es/                    # Elasticsearch 客户端
    ratelimit/             # Redis 限流中间件
    circuitbreaker/        # gobreaker 熔断器
    cache/                 # 本地缓存 + 热点检测
mianshiya-next-frontend/   # Next.js 前端
docker-compose.yml         # Docker Compose 编排
```

## 本地运行

```bash
# 后端
cd mianshiya-go-backend
go run ./cmd/server/

# 前端
cd mianshiya-next-frontend
npm install --legacy-peer-deps
npm run dev
```

## Docker 部署

```bash
docker compose up -d
# 访问 http://localhost
```

## 技术亮点

- **Go Clean Architecture**：Handler → Service → Repository 分层，依赖注入
- **分布式限流**：Redis + Lua 固定窗口，单机/分布式兼容
- **熔断降级**：gobreaker 保护 ES 调用，失败自动降级 MySQL LIKE
- **热点缓存**：Redis 计数器 + sync.Map 本地缓存，自动发现热 Key
- **ES 双写**：异步 goroutine 同步，不阻塞主流程
- **AI 重试**：题目生成带递增退避重试（1s→2s→3s）
- **跨平台部署**：多阶段 Docker 构建，Nginx 反向代理，一键部署
