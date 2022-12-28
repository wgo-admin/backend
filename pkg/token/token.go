package token

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	secret      string
	identityKey string
	roleKey     string
}

// ErrMissingHeader 表示 `Authorization` 请求头为空.
var ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")

var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey", "roleKey"}
	once   sync.Once
)

// 初始化配置
func Init(secret, identityKey, roleKey string) {
	once.Do(func() {
		if secret != "" {
			config.secret = secret
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
		if roleKey != "" {
			config.roleKey = roleKey
		}
	})
}

func Parse(tokenString string, key string) (string, string, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保 token 加密算法是预期的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(key), nil
	})
	// 解析失败
	if err != nil {
		return "", "", err
	}

	var identityKey string
	var roleKey string
	// 如果解析成功，从 token 中取出 token 的主题
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.identityKey].(string)
		roleKey = claims[config.roleKey].(string)
	}

	return identityKey, roleKey, nil
}

// ParseRequest 从请求头中获取令牌，并将其传递给 Parse 函数以解析令牌.
func ParseRequest(c *gin.Context) (string, string, error) {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return "", "", ErrMissingHeader
	}

	var t string
	// 从请求头中取出 token
	fmt.Sscanf(header, "Bearer %s", &t)

	return Parse(t, config.secret)
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(identityKey, roleKey string) (tokenString string, err error) {
	// Token 的内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		config.roleKey:     roleKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Add(100000 * time.Hour).Unix(),
	})

	// 签发 token
	return token.SignedString([]byte(config.secret))
}
