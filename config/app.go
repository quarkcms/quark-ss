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

	// 令牌加密key，默认自动生成，如果设置绝对不可泄漏
	"key": env.Get("APP_KEY"),
}
