package actions

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
)

type EditLink struct {
	actions.Link
}

// 初始化
func (p *EditLink) Init(name string) *EditLink {
	// 初始化父结构
	p.ParentInit()

	// 设置按钮类型,primary | ghost | dashed | link | text | default
	p.Type = "link"

	// 设置按钮大小,large | middle | small | default
	p.Size = "small"

	// 文字
	p.Name = name

	// 设置展示位置
	p.SetOnlyOnIndexTableRow(true)

	return p
}

// 跳转链接
func (p *EditLink) GetHref(c *fiber.Ctx) string {
	return "#/index?api=" + strings.Replace(strings.Replace(c.Path(), "/api/", "", -1), "/index", "/edit&id=${id}", -1)
}
