package session

import (
	fiber "github.com/gofiber/fiber/v2"
	session "github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New()
var ctx *fiber.Ctx

// 初始化
func Init(c *fiber.Ctx) {
	ctx = c
}

// 设置值
func Set(key string, value interface{}) error {

	sess, err := store.Get(ctx)
	if err != nil {
		panic(err)
	}
	sess.Set(key, value)

	return sess.Save()
}

// 获取值
func Get(key string) interface{} {

	sess, err := store.Get(ctx)
	if err != nil {
		panic(err)
	}
	result := sess.Get(key)

	return result
}
