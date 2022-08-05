package msg

import (
	"github.com/gofiber/fiber/v2"
)

var ctx *fiber.Ctx

const DEFAULT_MSG string = ""
const DEFAULT_URL string = ""
const DEFAULT_DATA string = ""

// 初始化
func Init(c *fiber.Ctx) {
	ctx = c
}

// 返回错误信息
func Error(msg string, url string) error {
	return ctx.JSON(fiber.Map{
		"component": "message",
		"status":    "error",
		"msg":       msg,
		"url":       url,
	})
}

// 返回正确信息
func Success(msg string, url string, data interface{}) error {
	return ctx.JSON(fiber.Map{
		"component": "message",
		"status":    "success",
		"msg":       msg,
		"url":       url,
		"data":      data,
	})
}
