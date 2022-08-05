package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/quarkcms/quark-go/database"
)

type Migrate struct {
	Command
}

// 初始化
func (p *Migrate) Init() *Migrate {
	p.Signature = "migrate"
	p.Description = ""

	return p
}

// 执行命令
func (p *Migrate) Handle() {

	// 执行迁移
	(&database.Migrate{}).Handle()

	color.Set(color.FgGreen)
	fmt.Println("Migration table created successfully.")
	color.Unset()
}
