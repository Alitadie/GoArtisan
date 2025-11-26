package main

import (
	"go-artisan/internal/bootstrap"

	"go.uber.org/fx"
)

// main 是程序的唯一入口
// 我们使用 Uber Fx 来管理整个应用程序的生命周期（依赖注入 + 启动/关闭钩子）
func main() {
	fx.New(
		// 1. 引入核心模块（配置、日志、数据库、路由、HTTPServer）
		bootstrap.Module,

		// 2. 这里的 Invoke 触发核心的启动逻辑
		// 只要我们在 bootstrap.Start 里写了 onStart 钩子，它就会在这里运行
		fx.Invoke(bootstrap.Start),
	).Run()
}
