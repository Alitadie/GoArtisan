package commands

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// 模板更新：自动引入 response 包和日志包，遵循依赖注入规范
const handlerTemplate = `package handler

import (
	"go-artisan/pkg/response"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type {{.Name}}Handler struct {
	logger *slog.Logger
	// 这里可以添加 service 依赖，例如: svc *service.{{.Name}}Service
}

// New{{.Name}}Handler 构造函数
func New{{.Name}}Handler(logger *slog.Logger) *{{.Name}}Handler {
	return &{{.Name}}Handler{
		logger: logger,
	}
}

// Index 示例方法
func (h *{{.Name}}Handler) Index(c *gin.Context) {
	// 示例：使用统一响应
	h.logger.Info("Accessing {{.Name}} Index")
	response.Success(c, gin.H{"module": "{{.Name}}", "action": "index"})
}
`

// NewMakeControllerCommand 保持不变... (省略部分并未修改逻辑)
func NewMakeControllerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:controller [name]",
		Short: "Create a new controller handler",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			titleName := strings.ToUpper(name[:1]) + name[1:]

			// 这里简单的转一下 snake_case，实际项目可以用 xstrings 库处理更复杂情况
			dirPath := "internal/http/handler"
			fileName := fmt.Sprintf("%s/%s_handler.go", dirPath, strings.ToLower(name))

			if err := os.MkdirAll(dirPath, 0755); err != nil {
				fmt.Printf("❌ Failed to create directory: %v\n", err)
				os.Exit(1)
			}

			if _, err := os.Stat(fileName); err == nil {
				fmt.Printf("❌ File already exists: %s\n", fileName)
				os.Exit(1)
			}

			f, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("❌ Failed to create file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()

			t := template.Must(template.New("handler").Parse(handlerTemplate))
			data := struct{ Name string }{Name: titleName}

			if err := t.Execute(f, data); err != nil {
				fmt.Printf("❌ Failed to execute template: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Controller created successfully: %s\n", fileName)
			fmt.Printf("👉 Don't forget to register it in internal/bootstrap/app.go and router.go!\n")
		},
	}
}
