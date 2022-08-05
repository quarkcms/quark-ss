package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
)

type Dashboard struct{}

// 列表页
func (p *Dashboard) Handle(c *fiber.Ctx) error {

	dashboard := &requests.Dashboard{}

	// 获取实例
	dashboardInstance := dashboard.Resource(c)

	// 断言DashboardComponentRender方法
	dashboardComponent := dashboardInstance.(interface {
		DashboardComponentRender(*fiber.Ctx, interface{}) interface{}
	}).DashboardComponentRender(c, dashboardInstance)

	// 断言Render方法
	component := dashboardInstance.(interface {
		Render(*fiber.Ctx, interface{}, interface{}) interface{}
	}).Render(c, dashboardInstance, dashboardComponent)

	return c.JSON(component)
}
