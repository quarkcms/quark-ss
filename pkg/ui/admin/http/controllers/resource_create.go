package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
)

type ResourceCreate struct{}

// 执行行为
func (p *ResourceCreate) Handle(c *fiber.Ctx) error {
	resourceCreate := &requests.ResourceCreate{}

	// 资源实例
	resourceInstance := resourceCreate.Resource(c)

	if resourceInstance == nil {
		return c.SendStatus(404)
	}

	// 断言BeforeCreating方法，获取初始数据
	data := resourceInstance.(interface {
		BeforeCreating(*fiber.Ctx) map[string]interface{}
	}).BeforeCreating(c)

	// 断言CreationComponentRender方法
	creationComponent := resourceInstance.(interface {
		CreationComponentRender(*fiber.Ctx, interface{}, map[string]interface{}) interface{}
	}).CreationComponentRender(c, resourceInstance, data)

	// 断言Render方法
	component := resourceInstance.(interface {
		Render(*fiber.Ctx, interface{}, interface{}) interface{}
	}).Render(c, resourceInstance, creationComponent)

	return c.JSON(component)
}
