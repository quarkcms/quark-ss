package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
)

type ResourceEditable struct{}

// 列表行内编辑
func (p *ResourceEditable) Handle(c *fiber.Ctx) error {
	resourceEditable := (&requests.ResourceEditable{}).HandleEditable(c)

	if resourceEditable == nil {
		return msg.Success("操作成功！", "", "")
	} else {
		return msg.Error("操作失败！", "")
	}
}
