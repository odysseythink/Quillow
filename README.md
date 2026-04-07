# Quillow

Quillow 是一个 AI 驱动的自托管个人财务管理系统，基于 Go + React/TypeScript 构建。

> 名字来源：Quill（羽毛笔，古典记账意象）+ Pillow（安心感）— 细致记账，安心理财。

## 功能特性

### 核心财务管理
- 多类型账户管理（资产、支出、收入、负债、现金）
- 交易记录（支出、收入、转账）
- 预算管理与预算限额
- 账单/订阅追踪
- 分类与标签体系
- 存钱罐（储蓄目标）
- 自动化规则引擎
- 定期交易
- 多币种支持与汇率

### AI 能力
- **智能分类** — 新交易自动识别分类，3 级 fallback（用户规则 → 本地模式匹配 → LLM）
- **自然语言记账** — 在聊天气泡中输入"昨天午饭35元"，AI 自动解析并创建交易
- **财务洞察** — 对话式问答（"这个月餐饮花了多少？""比上个月多花了多少？"）
- **自动学习** — 用户每次确认/修改分类，系统自动记住偏好

### 数据导入
- 微信支付账单 CSV 导入
- 支付宝账单 CSV 导入
- 自动格式检测与解析
- 导入前预览与去重

### 系统特性
- 多语言支持（中文、English、Deutsch）
- JWT 认证
- RESTful API
- 响应式前端（Ant Design）
- 自托管，数据完全自主

## 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go, Gin, GORM |
| 前端 | React 18, TypeScript, Ant Design, Redux Toolkit |
| 数据库 | MySQL / PostgreSQL / SQLite |
| AI | Claude API / OpenAI API + 本地模式匹配 |
| 认证 | JWT (golang-jwt) |
| 配置 | Viper + .env |
| 构建 | Makefile, Vite |

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+ / PostgreSQL 14+ / SQLite

### 安装与运行

```bash
# 克隆仓库
git clone https://github.com/odysseythink/Quillow.git
cd Quillow

# 配置环境变量
cp .env.example .env
# 编辑 .env 填入数据库连接信息

# 构建
make build

# 运行
make run
```

访问 http://localhost:8080

默认管理员账户：
- 邮箱：`admin@firefly.local`
- 密码：`firefly`

### 开发模式

```bash
make dev    # 同时启动前端 dev server 和后端
```

### 启用 AI 功能

编辑 `.env`：
```
AI_PROVIDER=claude    # 或 openai
AI_API_KEY=your-api-key
```

## 项目结构

```
cmd/server/          # 应用入口
internal/
├── entity/          # 领域实体
├── port/            # 接口定义
├── usecase/         # 业务逻辑
└── adapter/
    ├── handler/     # HTTP 处理器 (Gin)
    ├── repository/  # 数据访问 (GORM)
    └── transformer/ # 数据转换
pkg/
├── ai/              # AI 服务（分类、NLP、洞察）
├── importer/        # CSV 导入（微信、支付宝）
├── config/          # 配置管理
├── database/        # 数据库连接与迁移
├── jwt/             # JWT 认证
└── i18n/            # 国际化
web/                 # React 前端
```

## API 概览

| 类别 | 端点 |
|---|---|
| 认证 | `POST /api/v1/auth/login`, `/auth/refresh` |
| 账户 | `GET/POST/PUT/DELETE /api/v1/accounts` |
| 交易 | `GET/POST/DELETE /api/v1/transactions` |
| 预算 | `GET/POST/PUT/DELETE /api/v1/budgets` |
| AI 聊天 | `POST /api/v1/ai/chat` |
| AI 分类 | `POST /api/v1/ai/suggest` |
| 导入 | `POST /api/v1/import/preview`, `/import/confirm` |

## 许可证

AGPL-3.0
