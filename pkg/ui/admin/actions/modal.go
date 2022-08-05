package actions

import "github.com/gofiber/fiber/v2"

type Modal struct {
	Action
	Width int `json:"width"`
}

// 初始化
func (p *Modal) ParentInit() interface{} {
	p.ActionType = "modal"
	p.Width = 520

	return p
}

// 宽度
func (p *Modal) GetWidth() int {
	return p.Width
}

// 内容
func (p *Modal) GetBody(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	return nil
}

// 弹窗行为
func (p *Modal) GetActions(c *fiber.Ctx, resourceInstance interface{}) []interface{} {
	return []interface{}{}
}
