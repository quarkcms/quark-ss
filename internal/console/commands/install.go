package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/quarkcms/quark-go/database"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
)

type Install struct {
	Command
}

// 初始化
func (p *Install) Init() *Install {
	p.Signature = "install"
	p.Description = ""

	return p
}

// 执行命令
func (p *Install) Handle() {

	// 如果锁定文件存在则不执行安装步骤
	if utils.PathExist("install.lock") {
		color.Set(color.FgRed)
		fmt.Println("Install failed : The lock file exists in the root path, please delete it!")
		color.Unset()

		return
	}

	// 创建软连接
	storagePath := filepath.Join("..", "storage", "app", "public")
	SymlinkPath := filepath.Join("public", "storage")

	err := os.Symlink(storagePath, SymlinkPath)
	if err != nil {
		fmt.Print(err)
	}

	// 执行迁移
	(&database.Migrate{}).Handle()

	// 数据填充
	(&database.Seed{}).Handle()

	// 创建锁定文件
	file, _ := os.Create("install.lock")
	file.Close()

	color.Set(color.FgGreen)
	fmt.Println("The application have been installed.")
	color.Unset()
}
