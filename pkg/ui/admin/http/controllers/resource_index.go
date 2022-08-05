package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
)

type ResourceIndex struct{}

// 列表页
func (p *ResourceIndex) Handle(c *fiber.Ctx) error {

	resourceIndex := &requests.ResourceIndex{}

	// 资源实例
	resourceInstance := resourceIndex.Resource(c)

	if resourceInstance == nil {
		return c.SendStatus(404)
	}

	// 查询数据
	data := resourceIndex.IndexQuery(c)

	// 断言IndexComponentRender方法
	indexComponent := resourceInstance.(interface {
		IndexComponentRender(*fiber.Ctx, interface{}, interface{}) interface{}
	}).IndexComponentRender(c, resourceInstance, data)

	// 断言Render方法
	component := resourceInstance.(interface {
		Render(*fiber.Ctx, interface{}, interface{}) interface{}
	}).Render(c, resourceInstance, indexComponent)

	return c.JSON(component)
}
