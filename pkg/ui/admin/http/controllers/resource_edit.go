package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
)

type ResourceEdit struct{}

// 执行行为
func (p *ResourceEdit) Handle(c *fiber.Ctx) error {
	resourceEdit := &requests.ResourceEdit{}

	// 资源实例
	resourceInstance := resourceEdit.Resource(c)

	if resourceInstance == nil {
		return c.SendStatus(404)
	}

	data := resourceEdit.FillData(c)

	// 断言BeforeEditing方法，获取初始数据
	data = resourceInstance.(interface {
		BeforeEditing(*fiber.Ctx, map[string]interface{}) map[string]interface{}
	}).BeforeEditing(c, data)

	// 断言UpdateComponentRender方法
	updateComponent := resourceInstance.(interface {
		UpdateComponentRender(*fiber.Ctx, interface{}, map[string]interface{}) interface{}
	}).UpdateComponentRender(c, resourceInstance, data)

	// 断言Render方法
	component := resourceInstance.(interface {
		Render(*fiber.Ctx, interface{}, interface{}) interface{}
	}).Render(c, resourceInstance, updateComponent)

	return c.JSON(component)
}

// 获取表单初始化数据
func (p *ResourceEdit) Values(c *fiber.Ctx) error {
	resourceEdit := &requests.ResourceEdit{}

	// 资源实例
	resourceInstance := resourceEdit.Resource(c)

	if resourceInstance == nil {
		return c.SendStatus(404)
	}

	data := resourceEdit.FillData(c)

	// 断言BeforeEditing方法，获取初始数据
	data = resourceInstance.(interface {
		BeforeEditing(*fiber.Ctx, map[string]interface{}) map[string]interface{}
	}).BeforeEditing(c, data)

	return msg.Success("获取成功", "", data)
}
