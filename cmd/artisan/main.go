package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// 简单的 Handler 模板
// 使用 text/template 动态生成内容
const handlerTemplate = `package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type {{.Name}}Handler struct {}

func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{}
}

func (h *{{.Name}}Handler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from {{.Name}}"})
}
`

func main() {
	var rootCmd = &cobra.Command{
		Use:   "artisan",
		Short: "GoArtisan CLI Tool",
		Long:  "Helper utility to generate code and manage the GoArtisan application",
	}

	var makeControllerCmd = &cobra.Command{
		Use:   "make:controller [name]",
		Short: "Create a new controller handler",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// 1. 处理输入参数
			name := args[0]
			// 简单的首字母大写处理
			titleName := strings.ToUpper(name[:1]) + name[1:]

			// 2. 准备文件路径
			dirPath := "internal/http/handler"
			fileName := fmt.Sprintf("%s/%s_handler.go", dirPath, strings.ToLower(name))

			// 确保目录存在
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				fmt.Printf("❌ Failed to create directory: %v\n", err)
				os.Exit(1)
			}

			// 3. 创建文件
			f, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("❌ Failed to create file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()

			// 4. 解析模板
			t := template.Must(template.New("handler").Parse(handlerTemplate))
			data := struct{ Name string }{Name: titleName}

			if err := t.Execute(f, data); err != nil {
				fmt.Printf("❌ Failed to execute template: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Controller created successfully: %s\n", fileName)
		},
	}

	rootCmd.AddCommand(makeControllerCmd)

	// --- 修复点：捕获 Execute 的错误 ---
	if err := rootCmd.Execute(); err != nil {
		// Cobra 默认会打印错误信息，所以我们只需要确保以非零状态退出
		os.Exit(1)
	}
}
