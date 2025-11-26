package commands

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// 模板定义移到这里
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

// NewMakeControllerCommand 构造函数，返回 Command 指针
func NewMakeControllerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:controller [name]",
		Short: "Create a new controller handler",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			titleName := strings.ToUpper(name[:1]) + name[1:]

			dirPath := "internal/http/handler"
			fileName := fmt.Sprintf("%s/%s_handler.go", dirPath, strings.ToLower(name))

			if err := os.MkdirAll(dirPath, 0755); err != nil {
				fmt.Printf("❌ Failed to create directory: %v\n", err)
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
		},
	}
}
