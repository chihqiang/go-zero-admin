package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/x/errors"
	"time"
)

// JWT_KEY_ID 定义在 context 中存储用户ID的key
const JWT_KEY_ID = "id"

// JWT_KEY_TOKEN_TYPE 定义在 claims 中存储 Token 类型的 key
const JWT_KEY_TOKEN_TYPE = "token_type"

// TokenTypeAccess Access Token 类型
const TokenTypeAccess = "access"

// TokenTypeRefresh Refresh Token 类型
const TokenTypeRefresh = "refresh"

// JWTConfig JWT 配置
type JWTConfig struct {
	AccessSecret  string // Token 签名密钥（Access 和 Refresh 使用同一个密钥）
	AccessExpire  int64  `json:",default=7200"`              // Access Token 过期时间（秒），默认 2 小时
	RefreshExpire int64  `json:",default=604800"`            // Refresh Token 过期时间（秒），默认 7 天
	Iss           string `json:",default=go-zero-admin-api"` // 签发者
}

// JWTS JWT 签名工具
type JWTS struct {
	cfg *JWTConfig
}

// NewJWTS 创建 JWT 签名工具
//
// 注意：cfg 通常通过 go-zero 的 LoadConfig 加载，default 标签会自动填充默认值
func NewJWTS(cfg *JWTConfig) *JWTS {
	return &JWTS{cfg: cfg}
}

// GenerateAccessToken 生成 Access Token（短期有效）
//
// 参数：
//   - id: 用户ID
//
// 返回：
//   - token: 生成的 Access Token 字符串
//   - err: 错误信息
func (a *JWTS) GenerateAccessToken(id int64) (string, error) {
	// 参数校验
	if id <= 0 {
		return "", errors.New(400, "用户ID无效")
	}

	// 计算过期时间
	expireTime := time.Now().Add(time.Duration(a.cfg.AccessExpire) * time.Second)

	// 构建JWT载荷信息
	claims := jwt.MapClaims{
		JWT_KEY_ID:         id,                             // 用户ID
		JWT_KEY_TOKEN_TYPE: TokenTypeAccess,                // Token 类型
		"exp":              jwt.NewNumericDate(expireTime), // 过期时间
		"iat":              jwt.NewNumericDate(time.Now()), // 签发时间
		"iss":              a.cfg.Iss,                      // 签发者
		"nbf":              jwt.NewNumericDate(time.Now()), // 生效时间（Not Before）
	}

	// 使用HS256算法签名生成Token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.cfg.AccessSecret))
	if err != nil {
		return "", fmt.Errorf("生成Access Token失败: %w", err)
	}

	return token, nil
}

// GenerateRefreshToken 生成 Refresh Token（长期有效）
//
// 参数：
//   - id: 用户ID
//
// 返回：
//   - token: 生成的 Refresh Token 字符串
//   - err: 错误信息
func (a *JWTS) GenerateRefreshToken(id int64) (string, error) {
	// 参数校验
	if id <= 0 {
		return "", errors.New(400, "用户ID无效")
	}

	// 计算过期时间
	expireTime := time.Now().Add(time.Duration(a.cfg.RefreshExpire) * time.Second)

	// 构建JWT载荷信息
	claims := jwt.MapClaims{
		JWT_KEY_ID:         id,                             // 用户ID
		JWT_KEY_TOKEN_TYPE: TokenTypeRefresh,               // Token 类型
		"exp":              jwt.NewNumericDate(expireTime), // 过期时间
		"iat":              jwt.NewNumericDate(time.Now()), // 签发时间
		"iss":              a.cfg.Iss,                      // 签发者
		"nbf":              jwt.NewNumericDate(time.Now()), // 生效时间（Not Before）
	}

	// 使用HS256算法签名生成Token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.cfg.AccessSecret))
	if err != nil {
		return "", fmt.Errorf("生成Refresh Token失败: %w", err)
	}

	return token, nil
}

// GenerateTokenPair 生成 Access Token 和 Refresh Token 对
//
// 参数：
//   - id: 用户ID
//
// 返回：
//   - accessToken: Access Token
//   - refreshToken: Refresh Token
//   - err: 错误信息
func (a *JWTS) GenerateTokenPair(id int64) (accessToken string, refreshToken string, err error) {
	accessToken, err = a.GenerateAccessToken(id)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = a.GenerateRefreshToken(id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Generate 生成 JWT Token（兼容旧版本，内部调用 GenerateAccessToken）
//
// 参数：
//   - id: 用户ID
//
// 返回：
//   - token: 生成的 JWT Token 字符串
//   - err: 错误信息
func (a *JWTS) Generate(id int64) (string, error) {
	return a.GenerateAccessToken(id)
}

// GetID 从 context 中获取用户ID
//
// 参数：
//   - ctx: 上下文，通常从 gin.Context 或 http.Request.Context 获取
//
// 返回：
//   - id: 用户ID
//   - err: 错误信息
func (a *JWTS) GetID(ctx context.Context) (int64, error) {
	val := ctx.Value(JWT_KEY_ID)
	if val == nil {
		return 0, errors.New(401, "token 无效或已过期")
	}

	switch v := val.(type) {
	case int64:
		if v <= 0 {
			return 0, errors.New(401, "用户ID无效")
		}
		return v, nil
	case float64:
		if v <= 0 {
			return 0, errors.New(401, "用户ID无效")
		}
		return int64(v), nil
	case json.Number:
		id, err := v.Int64()
		if err != nil {
			return 0, errors.New(401, "用户ID格式错误")
		}
		if id <= 0 {
			return 0, errors.New(401, "用户ID无效")
		}
		return id, nil
	default:
		return 0, errors.New(401, fmt.Sprintf("不支持的用户ID类型: %T", val))
	}
}

// ParseAccessToken 解析 Access Token
//
// 参数：
//   - tokenString: Access Token 字符串
//
// 返回：
//   - claims: JWT 载荷
//   - err: 错误信息
func (a *JWTS) ParseAccessToken(tokenString string) (jwt.MapClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名算法错误: %v", token.Header["alg"])
		}
		return []byte(a.cfg.AccessSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析Access Token失败: %w", err)
	}

	// 验证Token是否有效
	if !token.Valid {
		return nil, errors.New(401, "Access Token无效")
	}

	// 获取载荷
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(401, "Access Token载荷格式错误")
	}

	// 验证Token类型
	if tokenType, ok := claims[JWT_KEY_TOKEN_TYPE]; ok {
		if tokenType != TokenTypeAccess {
			return nil, errors.New(401, "Token类型错误，期望Access Token")
		}
	}

	return claims, nil
}

// ParseRefreshToken 解析 Refresh Token
//
// 参数：
//   - tokenString: Refresh Token 字符串
//
// 返回：
//   - claims: JWT 载荷
//   - err: 错误信息
func (a *JWTS) ParseRefreshToken(tokenString string) (jwt.MapClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名算法错误: %v", token.Header["alg"])
		}
		return []byte(a.cfg.AccessSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析Refresh Token失败: %w", err)
	}

	// 验证Token是否有效
	if !token.Valid {
		return nil, errors.New(401, "Refresh Token无效")
	}

	// 获取载荷
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(401, "Refresh Token载荷格式错误")
	}

	// 验证Token类型
	if tokenType, ok := claims[JWT_KEY_TOKEN_TYPE]; ok {
		if tokenType != TokenTypeRefresh {
			return nil, errors.New(401, "Token类型错误，期望Refresh Token")
		}
	}

	return claims, nil
}

// ParseToken 解析 JWT Token（兼容旧版本，内部调用 ParseAccessToken）
//
// 参数：
//   - tokenString: JWT Token 字符串
//
// 返回：
//   - claims: JWT 载荷
//   - err: 错误信息
func (a *JWTS) ParseToken(tokenString string) (jwt.MapClaims, error) {
	return a.ParseAccessToken(tokenString)
}

// GetIDFromToken 从 Token 字符串中直接获取用户ID（可选功能）
//
// 参数：
//   - tokenString: JWT Token 字符串
//
// 返回：
//   - id: 用户ID
//   - err: 错误信息
func (a *JWTS) GetIDFromToken(tokenString string) (int64, error) {
	claims, err := a.ParseAccessToken(tokenString)
	if err != nil {
		return 0, err
	}

	// 获取用户ID
	val, ok := claims[JWT_KEY_ID]
	if !ok {
		return 0, errors.New(401, "Token中不包含用户ID")
	}

	// 类型断言
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case string:
		// 尝试将字符串转换为 int64
		var id int64
		if _, err := fmt.Sscanf(v, "%d", &id); err != nil {
			return 0, errors.New(401, "用户ID格式错误")
		}
		return id, nil
	default:
		return 0, errors.New(401, "用户ID类型错误")
	}
}

// RefreshAccessToken 使用 Refresh Token 刷新 Access Token
//
// 流程：
//  1. 验证 Refresh Token 是否有效
//  2. 从 Refresh Token 中提取用户ID
//  3. 生成新的 Access Token
//
// 参数：
//   - refreshTokenString: Refresh Token 字符串
//
// 返回：
//   - newAccessToken: 新生成的 Access Token
//   - err: 错误信息
//
// 示例：
//
//	newToken, err := jwtS.RefreshAccessToken(oldRefreshToken)
func (a *JWTS) RefreshAccessToken(refreshTokenString string) (newAccessToken string, err error) {
	// 解析 Refresh Token
	claims, err := a.ParseRefreshToken(refreshTokenString)
	if err != nil {
		return "", errors.New(401, "Refresh Token无效或已过期")
	}

	// 获取用户ID
	val, ok := claims[JWT_KEY_ID]
	if !ok {
		return "", errors.New(401, "Refresh Token中不包含用户ID")
	}

	// 类型断言，提取用户ID
	var userID int64
	switch v := val.(type) {
	case float64:
		userID = int64(v)
	case string:
		// 尝试将字符串转换为 int64
		if _, err := fmt.Sscanf(v, "%d", &userID); err != nil {
			return "", errors.New(401, "用户ID格式错误")
		}
	default:
		return "", errors.New(401, "用户ID类型错误")
	}

	if userID <= 0 {
		return "", errors.New(401, "用户ID无效")
	}

	// 生成新的 Access Token
	newAccessToken, err = a.GenerateAccessToken(userID)
	if err != nil {
		return "", fmt.Errorf("生成新的Access Token失败: %w", err)
	}

	return newAccessToken, nil
}

// RefreshTokenPair 刷新 Access Token 和 Refresh Token（双更新）
//
// 说明：
//   - 同时刷新 Access Token 和 Refresh Token（更安全）
//   - 旧的 Refresh Token 将失效
//
// 参数：
//   - refreshTokenString: 旧的 Refresh Token 字符串
//
// 返回：
//   - newAccessToken: 新生成的 Access Token
//   - newRefreshToken: 新生成的 Refresh Token
//   - err: 错误信息
func (a *JWTS) RefreshTokenPair(refreshTokenString string) (newAccessToken, newRefreshToken string, err error) {
	// 解析 Refresh Token
	claims, err := a.ParseRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", errors.New(401, "Refresh Token无效或已过期")
	}

	// 获取用户ID
	val, ok := claims[JWT_KEY_ID]
	if !ok {
		return "", "", errors.New(401, "Refresh Token中不包含用户ID")
	}

	// 类型断言，提取用户ID
	var userID int64
	switch v := val.(type) {
	case float64:
		userID = int64(v)
	case string:
		// 尝试将字符串转换为 int64
		if _, err := fmt.Sscanf(v, "%d", &userID); err != nil {
			return "", "", errors.New(401, "用户ID格式错误")
		}
	default:
		return "", "", errors.New(401, "用户ID类型错误")
	}

	if userID <= 0 {
		return "", "", errors.New(401, "用户ID无效")
	}

	// 生成新的 Token 对
	return a.GenerateTokenPair(userID)
}

// ValidateToken 验证 Token 是否有效（不提取用户信息，仅验证签名和过期时间）
//
// 参数：
//   - tokenString: JWT Token 字符串
//   - isRefreshToken: 是否为 Refresh Token（使用不同的密钥验证）
//
// 返回：
//   - valid: Token 是否有效
//   - err: 错误信息
func (a *JWTS) ValidateToken(tokenString string, isRefreshToken bool) (valid bool, err error) {
	if isRefreshToken {
		_, err = a.ParseRefreshToken(tokenString)
	} else {
		_, err = a.ParseAccessToken(tokenString)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
