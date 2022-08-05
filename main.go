package main

import (
	"embed"

	"github.com/quarkcms/quark-go/internal/console"
	"github.com/quarkcms/quark-go/internal/http"
)

//go:embed assets/*
var assets embed.FS

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

	// 控制台应用
	go consoleService()

	// 网站应用
	httpService()
}
