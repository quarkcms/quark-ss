package controllers

import (
	"bytes"

	"github.com/dchest/captcha"
	"github.com/go-basic/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/session"
)

type Captcha struct{}

// 创建session验证码
func (p *Captcha) Make(c *fiber.Ctx) error {

	digits := captcha.RandomDigits(4)

	session.Set("captcha", digits)

	image := captcha.NewImage(uuid.New(), digits, 110, 38)
	writer := bytes.Buffer{}
	image.WriteTo(&writer)

	return c.SendStream(bytes.NewReader(writer.Bytes()))
}

// 验证session验证码
func (p *Captcha) Check(digits string) bool {
	sessionCaptcha := session.Get("captcha")

	if digits == "" || sessionCaptcha == nil {
		return false
	}

	ns := make([]byte, len(digits))
	for i := range ns {
		d := digits[i]
		switch {
		case '0' <= d && d <= '9':
			ns[i] = d - '0'
		case d == ' ' || d == ',':
			// ignore
		default:
			return false
		}
	}

	return bytes.Equal(ns, sessionCaptcha.([]byte))
}

// 获取验证码ID
func (p *Captcha) GetID(c *fiber.Ctx) error {
	id := captcha.NewLen(4)

	return c.JSON(fiber.Map{
		"captcha_id":  id,
		"captcha_url": "/tool/captcha/" + id,
	})
}

// 创建ID验证码
func (p *Captcha) MakeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendString("Bad Request")
	}

	writer := bytes.Buffer{}
	captcha.WriteImage(&writer, id, 240, 80)

	return c.SendStream(bytes.NewReader(writer.Bytes()))
}

// 验证ID验证码
func (p *Captcha) CheckByID(id string, value string) bool {

	result := captcha.VerifyString(id, value)
	captcha.Reload(id)

	return result
}
