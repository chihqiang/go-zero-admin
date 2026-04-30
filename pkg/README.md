## 📖 功能模块

### 1. auth - JWT 认证

支持 Access Token + Refresh Token 双 Token 方案。

#### 配置

```go
// your-project.yaml
Auth:
  Secret: "your-secret-key"          # 签名密钥（Access 和 Refresh 共用）
  AccessExpire: 7200                 # Access Token 过期时间（秒），默认 2 小时
  RefreshExpire: 604800              # Refresh Token 过期时间（秒），默认 7 天
  Iss: "go-zero-api"                 # 签发者，默认 go-zero-api
```

#### 初始化

```go
import "your-project/pkg/auth"

cfg := &auth.JWTConfig{
    Secret:        "your-secret-key",
    AccessExpire:  7200,
    RefreshExpire: 604800,
    Iss:           "go-zero-api",
}

jwtS := auth.NewJWTS(cfg)
```

#### 生成 Token

```go
// 生成 Access Token
accessToken, err := jwtS.GenerateAccessToken(123)

// 生成 Refresh Token
refreshToken, err := jwtS.GenerateRefreshToken(123)

// 同时生成一对 Token（推荐）
accessToken, refreshToken, err := jwtS.GenerateTokenPair(123)

// 兼容旧版本（生成 Access Token）
token, err := jwtS.Generate(123)
```

#### 解析 Token

```go
// 解析 Access Token
claims, err := jwtS.ParseAccessToken(tokenString)

// 解析 Refresh Token
claims, err := jwtS.ParseRefreshToken(tokenString)

// 兼容旧版本
claims, err := jwtS.ParseToken(tokenString)
```

#### 刷新 Token

```go
// 使用 Refresh Token 刷新 Access Token
newAccessToken, err := jwtS.RefreshAccessToken(oldRefreshToken)

// 同时刷新 Access Token 和 Refresh Token（更安全）
newAccessToken, newRefreshToken, err := jwtS.RefreshTokenPair(oldRefreshToken)
```

#### 获取用户 ID

```go
// 从 context 中获取（配合中间件使用）
userID, err := jwtS.GetID(ctx)

// 直接从 Token 字符串获取
userID, err := jwtS.GetIDFromToken(tokenString)
```

#### 验证 Token

```go
// 验证 Token 是否有效
valid, err := jwtS.ValidateToken(tokenString, false) // false = Access Token
valid, err := jwtS.ValidateToken(tokenString, true)  // true = Refresh Token
```

---

### 2. orm - 数据库操作

基于 GORM 的数据库操作封装，支持 MySQL、PostgreSQL、SQLite。

#### 配置

```go
// your-project.yaml
Database:
  DBType: "mysql"                              # 数据库类型: mysql/postgres/sqlite
  DSN: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  MaxIdleConns: 10                              # 最大空闲连接数，默认 10
  MaxOpenConns: 100                             # 最大打开连接数，默认 100
  ConnMaxLifetime: 3600                         # 连接最长生命周期（秒），默认 3600
```

#### 使用方法

```go
import "your-project/pkg/orm"

// 打开数据库连接（失败返回 error）
db, err := orm.Open(&orm.Config{
    DBType: "mysql",
    DSN:    "root:root@tcp(localhost:3306)/test",
})

// 打开数据库连接（失败 panic，适合初始化阶段）
db := orm.MustOpen(&orm.Config{
    DBType: "mysql",
    DSN:    "root:root@tcp(localhost:3306)/test",
})
```

#### 支持的数据库

```go
// MySQL
db, _ := orm.Open(&orm.Config{DBType: "mysql", DSN: "root:root@tcp(localhost:3306)/test"})

// PostgreSQL
db, _ := orm.Open(&orm.Config{
    DBType: "postgres",
    DSN:    "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable",
})

// SQLite
db, _ := orm.Open(&orm.Config{DBType: "sqlite", DSN: "./test.db"})
```

---

### 3. hash - 密码哈希

```go
import "your-project/pkg/hash"

// 生成密码哈希
hashStr, err := hash.BcryptHash("password")

// 验证密码
match := hash.BcryptCheck("password", hashStr)

// 生成 MD5（用于第三方接口签名等场景）
md5Str := hash.MD5("hello")
```

---

## 🔧 工具脚本

### init-env.sh - 环境初始化

在项目根目录执行，自动检查和安装依赖：

```bash
bash pkg/init-env.sh
```

功能：
- 检查 Go 版本
- 安装 goctl 工具
- 安装 go-zero 依赖
- 初始化 go.mod（如不存在）

---

## 📝 注意事项

1. **JWT Secret 必须保密**：不要提交到代码仓库，建议使用环境变量或配置文件加密
2. **Refresh Token 安全**：Refresh Token 有效期较长，建议存储在安全的地方（如 HttpOnly Cookie）
3. **数据库配置**：生产环境请务必修改默认的连接池参数

---