package middleware

import (
	"go-artisan/pkg/version"

	"github.com/gin-gonic/gin"
)

func VersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 添加 HTTP 响应头
		c.Header("X-App-Version", version.GitTag)
		c.Header("X-App-Commit", version.GitCommit)
		c.Next()
	}
}
