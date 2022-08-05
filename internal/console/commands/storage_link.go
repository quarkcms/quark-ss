package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type StorageLink struct {
	Command
}

// 初始化
func (p *StorageLink) Init() *StorageLink {
	p.Signature = "storage:link"
	p.Description = ""

	return p
}

// 执行命令
func (p *StorageLink) Handle() {

	// 创建软连接
	storagePath := filepath.Join("..", "storage", "app", "public")
	SymlinkPath := filepath.Join("public", "storage")

	err := os.Symlink(storagePath, SymlinkPath)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		return
	}

	color.Set(color.FgGreen)
	fmt.Println("The [public\\storage] link has been connected to [storage\\app\\public].")
	fmt.Println("The links have been created.")
	color.Unset()
}
