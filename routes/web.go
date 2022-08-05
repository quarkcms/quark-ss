package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/http/controllers"
)

type Web struct{}

// web路由
func (p *Web) Route(app *fiber.App) {

	// hello world
	app.Get("/", (&controllers.Index{}).Index)
}
