package actions

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
)

type CreateLink struct {
	actions.Link
}

// 初始化
func (p *CreateLink) Init(name string) *CreateLink {
	// 初始化父结构
	p.ParentInit()

	// 类型
	p.Type = "primary"

	// 图标
	p.Icon = "plus-circle"

	// 文字
	p.Name = "创建" + name

	// 设置展示位置
	p.SetOnlyOnIndex(true)

	return p
}

// 跳转链接
func (p *CreateLink) GetHref(c *fiber.Ctx) string {
	return "#/index?api=" + strings.Replace(strings.Replace(c.Path(), "/api/", "", -1), "/index", "/create", -1)
}
