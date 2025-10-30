# User Management System - Go + MySQL

一个基于 Go 和 MySQL 的企业级用户管理系统，包含用户认证和文章管理功能。

## 功能特性

### 用户管理
- ✅ 用户注册
- ✅ 用户登录（JWT认证）
- ✅ 获取当前用户信息
- ✅ 更新用户信息
- ✅ 删除用户账户

### 文章管理
- ✅ 创建文章（需认证）
- ✅ 查看文章列表（公开访问）
- ✅ 查看文章详情（公开访问）
- ✅ 查看我的文章列表（需认证）
- ✅ 更新文章（仅作者）
- ✅ 删除文章（仅作者）

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT (golang-jwt/jwt)
- **密码加密**: bcrypt
- **配置管理**: Viper
- **数据验证**: validator/v10
- **数据库迁移**: gormigrate

## 项目结构

```
user-management/
├── cmd/                  # 应用入口点
│   └── api/
│       └── main.go       # 主程序入口
├── internal/             # 私有应用代码
│   ├── api/              # 控制器层
│   │   ├── handlers/     # 路由处理器
│   │   ├── middleware/   # 中间件
│   │   └── routes.go     # 路由定义
│   ├── config/           # 配置管理
│   ├── domain/           # 领域模型
│   ├── repository/       # 数据访问层
│   ├── service/          # 业务逻辑层
│   └── util/             # 工具函数
├── pkg/                  # 可公开重用的包
│   ├── logger/           # 日志工具
│   └── security/         # 安全相关（JWT、密码加密）
├── migrations/           # 数据库迁移
├── configs/              # 配置文件
│   └── config.yaml       # 默认配置
├── docs/                 # 文档
│   └── CONFIG.md         # 配置说明文档
├── .env.example          # 环境变量示例
└── README.md             # 项目说明
```

## 快速开始

> **💡 配置说明**：详细的配置管理说明请查看 [docs/CONFIG.md](docs/CONFIG.md)

### 1. 环境要求

- Go 1.24.4 或更高版本
- MySQL 5.7 或更高版本

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

创建 MySQL 数据库：

```sql
CREATE DATABASE user_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 配置文件

项目支持多种配置方式，推荐使用 `.env` 文件（开发环境）：

**方式 1：使用 .env 文件（推荐开发环境）**：

```bash
# 1. 复制配置文件
cp .env.example .env

# 2. 修改 .env 文件
vim .env

# 根据您的 MySQL 配置选择：
# - 如果 MySQL 无密码（本地开发）：DB_PASSWORD=
# - 如果 MySQL 有密码：DB_PASSWORD=your_password

# 3. 运行项目（会自动加载 .env）
go run cmd/api/main.go
```

**快速开始示例**：

```bash
# 无密码 MySQL（常见于本地开发）
cp .env.nopassword.example .env

# 有密码 MySQL
cp .env.withpassword.example .env
vim .env  # 修改 DB_PASSWORD 为实际密码
```

**方式 2：直接修改 config.yaml**：

```yaml
# configs/config.yaml
database:
  username: "root"
  password: ""  # 空密码，或填入实际密码

jwt:
  secretKey: "dev-secret-key"
```

**方式 3：系统环境变量（生产环境）**：

```bash
export DB_PASSWORD=your_password
export JWT_SECRET_KEY=your-secret-key
go run cmd/api/main.go
```

> 📖 配置优先级：`.env 文件` > `configs/config.yaml` > `默认值`
>
> ⚠️ **重要**：如果 `.env` 中设置 `DB_PASSWORD=`（空值），将使用无密码连接，不会使用 `config.yaml` 中的默认值
>
> 详细配置说明（包括 Docker、K8s 部署）请查看 [docs/CONFIG.md](docs/CONFIG.md)

### 5. 运行项目

```bash
# 方式1: 直接运行
go run cmd/api/main.go

# 方式2: 编译后运行
go build -o bin/api cmd/api/main.go
./bin/api
```

服务器将在 `http://localhost:8080` 启动。

## API 文档

### 公开接口

#### 1. 用户注册
```bash
POST /api/v1/users/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User"
}
```

#### 2. 用户登录
```bash
POST /api/v1/users/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}

# 响应
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 3. 获取文章列表
```bash
GET /api/v1/articles?page=1&page_size=10
```

#### 4. 获取文章详情
```bash
GET /api/v1/articles/:id
```

### 需要认证的接口

**所有需要认证的接口都需要在请求头中携带 JWT token：**
```
Authorization: Bearer <your_jwt_token>
```

#### 1. 获取当前用户信息
```bash
GET /api/v1/users/me
Authorization: Bearer <token>
```

#### 2. 更新用户信息
```bash
PUT /api/v1/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "full_name": "Updated Name",
  "email": "newemail@example.com"
}
```

#### 3. 删除用户账户
```bash
DELETE /api/v1/users/me
Authorization: Bearer <token>
```

#### 4. 创建文章
```bash
POST /api/v1/articles
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My First Article",
  "content": "This is the content of my article..."
}
```

#### 5. 获取我的文章列表
```bash
GET /api/v1/users/me/articles?page=1&page_size=10
Authorization: Bearer <token>
```

#### 6. 更新文章
```bash
PUT /api/v1/articles/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

#### 7. 删除文章
```bash
DELETE /api/v1/articles/:id
Authorization: Bearer <token>
```

## 测试示例

使用 curl 进行测试：

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# 2. 登录获取 token
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 3. 使用 token 创建文章（替换 <your_token>）
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "title": "My First Article",
    "content": "This is a great article about Go programming!"
  }'

# 4. 查看文章列表
curl http://localhost:8080/api/v1/articles
```

## 数据库迁移

项目使用 GORM 的自动迁移功能，首次运行时会自动创建表结构：
- `users` - 用户表
- `articles` - 文章表

如需重置数据库，可以删除数据库后重新创建：
```sql
DROP DATABASE user_management;
CREATE DATABASE user_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 开发建议

### 类比 Python/FastAPI/Django

如果你熟悉 Python 的 FastAPI 或 Django，这里是一些概念对应：

| Go (本项目) | Python (FastAPI/Django) |
|-------------|------------------------|
| handlers    | routes / views         |
| service     | service / business logic |
| repository  | ORM models / database layer |
| domain      | models / schemas       |
| middleware  | middleware             |
| Gin         | FastAPI / Django       |
| GORM        | SQLAlchemy / Django ORM |
| JWT         | python-jose / PyJWT    |
| Viper       | python-dotenv / settings |

### 安全注意事项

1. **生产环境**：
   - 修改 JWT secret key
   - 使用环境变量管��敏感配置
   - 启用 HTTPS
   - 配置适当的 CORS 策略

2. **密码安全**：
   - 已使用 bcrypt 加密
   - 密码最小长度 6 位（可在 domain/user.go 修改）

3. **数据验证**：
   - 使用 validator 进行输入验证
   - 前端也应做相应验证

## 下一步改进

- [ ] 添加单元测试
- [ ] 添加集成测试
- [ ] 实现刷新 token 机制
- [ ] 添加角色权限管理
- [ ] 添加日志中间件
- [ ] 添加 API 限流
- [ ] 添加 Swagger 文档
- [ ] 实现软删除恢复功能
- [ ] 添加文章分类和标签
- [ ] 实现文章搜索功能

## 许可证

MIT License

## 联系方式

如有问题，请提 Issue。
