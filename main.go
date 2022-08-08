package main

import (
	"embed"

	netHttp "net/http"
	_ "net/http/pprof"

	"github.com/quarkcms/quark-go/config"
	"github.com/quarkcms/quark-go/internal/console"
	"github.com/quarkcms/quark-go/internal/http"
)

//go:embed assets/*
var assets embed.FS

// 性能分析工具
func pprofService() {
	if config.App["pprof_server"].(string) == "true" {

		// 服务实例
		netHttp.ListenAndServe(config.App["pprof_host"].(string), nil)
	}
}

// 网站服务
func httpService() {

	// 服务实例
	kernel := &http.Kernel{}

	// 启动服务
	kernel.Run(assets)
}

// 控制台服务
func consoleService() {
	// 服务实例
	kernel := &console.Kernel{}

	// 启动服务
	kernel.Run(assets)
}

func main() {

	// 性能分析工具
	go pprofService()

	// 控制台应用
	go consoleService()

	// 网站应用
	httpService()
}
