package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务码 (0 或 200 代表成功，非 0 代表失败)
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据载荷
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200, // 或者 0，看公司规范
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应 (支持自定义 HTTP 状态码)
func Error(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus, // 这里简化处理，也可以传入具体业务错误码
		Message: msg,
		Data:    nil,
	})
}

// ValidationError 表单验证失败响应
// 模仿 Laravel: 422 Unprocessable Entity
func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Code:    422,
		Message: "Validation failed",
		Data:    errors,
	})
}
