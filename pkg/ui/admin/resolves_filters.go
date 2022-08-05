package admin

import "github.com/gofiber/fiber/v2"

//定义筛选表单
func (p *Resource) Filters(c *fiber.Ctx) []interface{} {
	return []interface{}{}
}
