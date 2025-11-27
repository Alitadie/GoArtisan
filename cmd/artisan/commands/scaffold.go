package commands

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/template"

	"go-artisan/internal/config"

	// 引入数据库驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

// Column 描述数据库列
type Column struct {
	Name string
	Type string
	Json string
}

// 转换 MySQL 类型到 Go 类型
func mysqlTypeToGo(mysqlType string) string {
	if strings.Contains(mysqlType, "int") {
		return "int" // 实际上要根据 unsigned 等区分，简化处理
	} else if strings.Contains(mysqlType, "datetime") || strings.Contains(mysqlType, "timestamp") {
		return "time.Time"
	}
	return "string"
}

// 生成用的数据包
type ScaffoldData struct {
	TableName   string
	StructName  string
	Columns     []Column
	PackageName string
}

const modelTemplate = `package domain

import "time"

// {{.StructName}} mapped from table {{.TableName}}
type {{.StructName}} struct {
{{- range .Columns }}
	{{ .Name }} {{ .Type }} ` + "`" + `json:"{{ .Json }}"` + "`" + `
{{- end }}
}
`

// NewMakeScaffoldCommand
// 使用: go run cmd/artisan/main.go make:scaffold users
func NewMakeScaffoldCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "make:scaffold [table_name]",
		Short: "Generate Domain/Model from existing Database Table",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tableName := args[0]
			fmt.Printf("🏗️  Scaffolding for table: %s...\n", tableName)

			// 1. 连接数据库 (读取列信息)
			db, err := sql.Open("mysql", cfg.Database.DSN)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			// 2. 查询表结构 (Information Schema 方式或者直接 Select Limit 0)
			rows, err := db.Query(fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = '%s' ORDER BY ordinal_position", tableName))
			if err != nil {
				fmt.Printf("❌ Failed to query schema: %v\n", err)
				os.Exit(1)
			}
			defer rows.Close()

			var columns []Column
			for rows.Next() {
				var colName, colType string
				if err := rows.Scan(&colName, &colType); err != nil {
					continue
				}

				// 简单的名字处理: user_email -> UserEmail
				goName := toTitle(colName)
				columns = append(columns, Column{
					Name: goName,
					Type: mysqlTypeToGo(colType),
					Json: colName, // json tag 保持下划线
				})
			}

			// 3. 准备数据
			data := ScaffoldData{
				TableName:  tableName,
				StructName: toTitle(tableName), // e.g. users -> Users (需处理单复数，简化处理)
				Columns:    columns,
			}

			// 4. 生成 Domain Model 文件
			// 实际项目你还需要生成 Service / Handler / Repo
			fileName := fmt.Sprintf("internal/domain/%s.go", strings.ToLower(tableName))
			generateFile(fileName, modelTemplate, data)

			fmt.Printf("✅ Model generated: %s\n", fileName)
		},
	}
}

// 辅助函数: 生成文件
func generateFile(path string, tmpl string, data interface{}) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t := template.Must(template.New("scaffold").Parse(tmpl))
	if err := t.Execute(f, data); err != nil {
		panic(err)
	}
}

// 辅助: 下划线转大驼峰 users_role -> UsersRole
func toTitle(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}
