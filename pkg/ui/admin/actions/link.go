package actions

import "github.com/gofiber/fiber/v2"

type Link struct {
	Action
	Href   string `json:"href"`
	Target string `json:"target"`
}

// 初始化
func (p *Link) ParentInit() interface{} {
	p.ActionType = "link"
	p.Target = "_self"

	return p
}

/**
 * 获取跳转链接
 *
 * @return string
 */
func (p *Link) GetHref(c *fiber.Ctx) string {
	return p.Href
}

/**
 * 相当于 a 链接的 target 属性，href 存在时生效
 *
 * @return string
 */
func (p *Link) GetTarget(c *fiber.Ctx) string {
	return p.Target
}
