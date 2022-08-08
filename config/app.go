package config

import (
	"github.com/quarkcms/quark-go/pkg/framework/env"
)

var App = map[string]interface{}{
	// 应用名称
	"name": env.Get("APP_NAME", "QuarkGo"),

	// 服务地址
	"host": env.Get("APP_HOST", "127.0.0.1:3000"),

	// 开发者模式
	"debug": env.Get("APP_DEBUG", "false"),

	// 性能分析工具，默认关闭
	"pprof_server": env.Get("APP_PPROF_SERVER", "false"),

	// 性能分析工具服务地址
	"pprof_host": env.Get("APP_PPROF_HOST", "127.0.0.1:8000"),

	// 令牌加密key，默认自动生成，如果设置绝对不可泄漏
	"key": env.Get("APP_KEY"),
}
