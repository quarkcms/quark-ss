package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type Index struct{}

// Index
func (p *Index) Index(c *fiber.Ctx) error {

	return c.SendString("Hello, World!")
}
