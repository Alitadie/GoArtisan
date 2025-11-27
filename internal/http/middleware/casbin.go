package middleware

import (
	"strconv"

	"go-artisan/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求的 URL 和 Method
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 2. 获取当前用户 (从 AuthMiddleware 设置的 Context 里取)
		uid, exists := c.Get(ContextUserIDKey) // 假设你在 auth.go 里 set 的 key 是这个
		if !exists {
			response.Error(c, 401, "Unauthenticated")
			c.Abort()
			return
		}

		// 假设我们的 Subject 是 "user:1", "user:2" 这种格式
		sub := "user:" + strconv.Itoa(int(uid.(uint)))

		// 3. 检查权限
		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			response.Error(c, 500, "Permission check error")
			c.Abort()
			return
		}

		if !ok {
			response.Error(c, 403, "You don't have permission")
			c.Abort()
			return
		}

		c.Next()
	}
}
