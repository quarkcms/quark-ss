package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
)

type ResourceDetail struct{}

// 执行行为
func (p *ResourceDetail) Handle(c *fiber.Ctx) error {
	ResourceDetail := &requests.ResourceDetail{}

	// 资源实例
	resourceInstance := ResourceDetail.Resource(c)

	if resourceInstance == nil {
		return c.SendStatus(404)
	}

	data := ResourceDetail.FillData(c)

	// 断言方法，获取初始数据
	data = resourceInstance.(interface {
		BeforeDetailShowing(*fiber.Ctx, map[string]interface{}) map[string]interface{}
	}).BeforeDetailShowing(c, data)

	// 断言方法
	detailComponent := resourceInstance.(interface {
		DetailComponentRender(*fiber.Ctx, interface{}, map[string]interface{}) interface{}
	}).DetailComponentRender(c, resourceInstance, data)

	// 断言Render方法
	component := resourceInstance.(interface {
		Render(*fiber.Ctx, interface{}, interface{}) interface{}
	}).Render(c, resourceInstance, detailComponent)

	return c.JSON(component)
}

// 获取表单初始化数据
func (p *ResourceDetail) Values(c *fiber.Ctx) error {
	ResourceDetail := &requests.ResourceDetail{}

	// 资源实例
	resourceInstance := ResourceDetail.Resource(c)

	if resourceInstance == nil {
		return c.SendStatus(404)
	}

	data := ResourceDetail.FillData(c)

	data = resourceInstance.(interface {
		BeforeDetailShowing(*fiber.Ctx, map[string]interface{}) map[string]interface{}
	}).BeforeDetailShowing(c, data)

	return msg.Success("获取成功", "", data)
}
