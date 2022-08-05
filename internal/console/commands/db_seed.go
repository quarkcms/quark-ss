package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/quarkcms/quark-go/database"
)

type DBSeed struct {
	Command
}

// 初始化
func (p *DBSeed) Init() *DBSeed {
	p.Signature = "db:seed"
	p.Description = ""

	return p
}

// 执行命令
func (p *DBSeed) Handle() {

	// 数据填充
	(&database.Seed{}).Handle()

	color.Set(color.FgGreen)
	fmt.Println("Database seeding completed successfully.")
	color.Unset()
}
