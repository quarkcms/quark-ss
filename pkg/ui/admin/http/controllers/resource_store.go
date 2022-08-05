package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
	"gorm.io/gorm"
)

type ResourceStore struct{}

// 执行行为
func (p *ResourceStore) Handle(c *fiber.Ctx) error {
	result := (&requests.ResourceStore{}).HandleStore(c)

	if value, ok := result.(error); ok {
		errorMsg := value.Error()
		if errorMsg != "" {
			return msg.Error(errorMsg, "")
		} else {
			return value
		}
	}

	if value, ok := result.(*gorm.DB); ok {
		if value.Error == nil {
			return msg.Success("操作成功！", "/index?api=admin/"+c.Params("resource")+"/index", "")
		} else {
			return msg.Error("操作失败！", "")
		}
	}

	return msg.Error("错误！", "")
}
