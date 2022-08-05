package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/quarkcms/quark-go/internal/shadowsocks"
)

type Shadowsocks struct {
	Command
}

// 初始化
func (p *Shadowsocks) Init() *Shadowsocks {
	p.Signature = "ss"
	p.Description = ""

	return p
}

// 执行命令
func (p *Shadowsocks) Handle() {

	shadowsocks.Start()

	color.Set(color.FgGreen)
	fmt.Println("The ss service started successfully.")
	color.Unset()
}

func StartShadowsocks() {
}
