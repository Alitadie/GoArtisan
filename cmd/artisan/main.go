package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

// 简单的 Handler 模板
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
	var rootCmd = &cobra.Command{Use: "artisan"}

	var makeControllerCmd = &cobra.Command{
		Use:   "make:controller [name]",
		Short: "Create a new controller handler",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fileName := fmt.Sprintf("internal/http/handler/%s_handler.go", name)

			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			t := template.Must(template.New("handler").Parse(handlerTemplate))
			data := struct{ Name string }{Name: name} // 注意大小写处理，这里简化了

			if err := t.Execute(f, data); err != nil {
				panic(err)
			}
			fmt.Printf("✅ Controller created: %s\n", fileName)
		},
	}

	rootCmd.AddCommand(makeControllerCmd)
	rootCmd.Execute()
}
