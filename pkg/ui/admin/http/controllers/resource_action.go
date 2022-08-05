package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
)

type ResourceAction struct{}

// 执行行为
func (p *ResourceAction) Handle(c *fiber.Ctx) error {
	return (&requests.ResourceAction{}).HandleAction(c)
}
