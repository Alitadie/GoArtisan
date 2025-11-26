package middleware

import (
	"strings"

	"go-artisan/pkg/auth"
	"go-artisan/pkg/response"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

// AuthMiddleware 接收配置中的 Secret
// 由于 Middleware 初始化在 Router 构造时，可以通过传参注入
func AuthMiddleware() gin.HandlerFunc {
	// TODO: 从配置中获取 Secret，目前写死必须与 Login 保持一致
	secret := "KeepItSecretKeepItSafe!GoArtisanKey"

	return func(c *gin.Context) {
		// 1. 获取 Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, "Authorization header required")
			c.Abort()
			return
		}

		// 2. 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, "Invalid authorization format")
			c.Abort()
			return
		}

		// 3. 校验 Token
		claims, err := auth.ParseToken(parts[1], secret)
		if err != nil {
			response.Error(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		// 4. 将 ID 注入上下文，后续 Controller 可以通过 c.Get("userID") 获取
		c.Set(ContextUserIDKey, claims.UserID)

		c.Next()
	}
}
