package version

import (
	"fmt"
	"runtime"
)

var (
	// 以下变量将在编译时通过 -ldflags 注入
	GitTag    = "dev"
	GitCommit = "none"
	BuildTime = "unknown"
)

// FullVersion 返回完整的版本信息字符串
func FullVersion() string {
	return fmt.Sprintf(
		"Version: %s\nCommit: %s\nBuilt: %s\nGo: %s\nOS/Arch: %s/%s",
		GitTag, GitCommit, BuildTime, runtime.Version(), runtime.GOOS, runtime.GOARCH,
	)
}

// Map 返回结构化信息，用于 API 响应
func Map() map[string]string {
	return map[string]string{
		"version":    GitTag,
		"commit":     GitCommit,
		"build_time": BuildTime,
	}
}
