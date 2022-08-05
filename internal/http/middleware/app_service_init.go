package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/framework/session"
)

type AppServiceInit struct{}

// 全局中间件
func (p *AppServiceInit) Handle(c *fiber.Ctx) error {

	// 初始化msg
	msg.Init(c)

	// 初始化session
	session.Init(c)

	return c.Next()
}
