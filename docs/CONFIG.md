# 配置管理说明

## 配置优先级（从高到低）

```
.env 文件 > configs/config.yaml > 代码默认值
```

> **重要**：项目启动时会自动加载项目根目录的 `.env` 文件（如果存在）

## 配置方式

### 方式 1：使用 .env 文件（推荐开发环境）

最简单的方式，类似 Python 的 `python-dotenv`：

**.env 文件**：
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_NAME=user_management

# JWT配置
JWT_SECRET_KEY=dev-secret-key

# 服务器配置（可选）
SERVER_PORT=8080
```

**快速开始**：
```bash
# 1. 复制示例文件
cp .env.example .env

# 2. 修改配置
vim .env

# 3. 运行（会自动加载 .env）
go run cmd/api/main.go
```

**优点**：
- ✅ 类似 Python/Node.js 开发习惯
- ✅ 所有配置集中在一个文件
- ✅ 不需要修改 config.yaml
- ✅ 敏感信息不进入版本控制（.env 在 .gitignore 中）

---

### 方式 2：仅使用 config.yaml（传统方式）

最简单的方式，直接修改 `configs/config.yaml`：

```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "your_password"  # 直接写在这里
  name: "user_management"

jwt:
  secretKey: "dev-secret-key"
  expirationHours: 24
```

启动应用：
```bash
go run cmd/api/main.go
```

---

### 方式 3：系统环境变量（生产环境）

基础配置放在 `config.yaml`，敏感信息用环境变量覆盖：

**configs/config.yaml**（非敏感配置）：
```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: ""  # 空着，用环境变量设置
  name: "user_management"

jwt:
  secretKey: ""  # 空着，用环境变量设置
  expirationHours: 24
```

**启动时设置环境变量**：
```bash
# Linux/Mac
export DB_PASSWORD=your_password
export JWT_SECRET_KEY=your-secret-key
go run cmd/api/main.go

# 或者一行命令
DB_PASSWORD=your_password JWT_SECRET_KEY=your-secret-key go run cmd/api/main.go
```

---

### 方式 3：仅使用环境变量（生产环境）

不需要 config.yaml，全部用环境变量：

```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_PORT=3306
export DB_USERNAME=root
export DB_PASSWORD=your_password
export DB_NAME=user_management
export JWT_SECRET_KEY=your-secret-key

go run cmd/api/main.go
```

---

## 支持的环境变量

| 环境变量 | 对应配置 | 说明 | 默认值 |
|---------|---------|------|--------|
| `SERVER_PORT` | server.port | 服务器端口 | 8080 |
| `DB_HOST` | database.host | 数据库主机 | localhost |
| `DB_PORT` | database.port | 数据库端口 | 3306 |
| `DB_USERNAME` | database.username | 数据库用户名 | root |
| `DB_PASSWORD` | database.password | 数据库密码 | - |
| `DB_NAME` | database.name | 数据库名称 | user_management |
| `JWT_SECRET_KEY` | jwt.secretKey | JWT密钥 | - |

---

## Docker / K8s 部署示例

### Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USERNAME=root
      - DB_PASSWORD=root123
      - DB_NAME=user_management
      - JWT_SECRET_KEY=super-secret-key-change-in-production
    depends_on:
      - mysql

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root123
      - MYSQL_DATABASE=user_management
    ports:
      - "3306:3306"
```

### Kubernetes ConfigMap + Secret

**configmap.yaml**（非敏感配置）：
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  SERVER_PORT: "8080"
  DB_HOST: "mysql-service"
  DB_PORT: "3306"
  DB_USERNAME: "root"
  DB_NAME: "user_management"
```

**secret.yaml**（敏感配置）：
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
stringData:
  DB_PASSWORD: "your_password"
  JWT_SECRET_KEY: "your-secret-key"
```

**deployment.yaml**：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-management
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: app
        image: user-management:latest
        envFrom:
        - configMapRef:
            name: app-config
        - secretRef:
            name: app-secret
```

---

## 常见问题

### Q: 我应该使用哪种配置方式？

**A:** 根据环境选择：
- **开发环境**：方式 1（仅 config.yaml），简单直接
- **测试环境**：方式 2（config.yaml + 环境变量），灵活性好
- **生产环境**：方式 2 或 3，敏感信息必须用环境变量

### Q: .env 文件在 Go 项目中必需吗？

**A:** 不必需。Go 项目通常直接使用环境变量，不像 Node.js/Python 那样依赖 .env 文件。

如果你习惯使用 .env 文件，可以安装 `godotenv` 库：
```bash
go get github.com/joho/godotenv
```

然后在 main.go 中加载：
```go
import "github.com/joho/godotenv"

func main() {
    godotenv.Load() // 加载 .env 文件
    // ...
}
```

### Q: 如何验证配置是否正确？

**A:** 在 main.go 中添加日志输出：
```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 验证配置（不输出敏感信息）
    log.Printf("Server Port: %s", cfg.Server.Port)
    log.Printf("Database Host: %s", cfg.Database.Host)
    log.Printf("Database Name: %s", cfg.Database.Name)
    // 不要输出密码和密钥！
}
```

---

## 最佳实践

1. **开发环境**：直接修改 `configs/config.yaml`
2. **生产环境**：敏感信息（密码、密钥）必须用环境变量
3. **版本控制**：
   - ✅ 提交 `configs/config.yaml`（不含敏感信息）
   - ✅ 提交 `.env.example`（示例）
   - ❌ 不要提交 `.env`（已在 .gitignore）
4. **安全性**：
   - 生产环境的 JWT 密钥必须足够复杂（至少32位）
   - 定期轮换密钥
   - 使用密钥管理服务（AWS Secrets Manager、HashiCorp Vault 等）
