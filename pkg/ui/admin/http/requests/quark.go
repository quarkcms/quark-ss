package requests

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin"
	"gorm.io/gorm"
)

type Quark struct{}

// 资源
func (p *Quark) Resource(c *fiber.Ctx) interface{} {
	var resourceInstance interface{}
	for _, provider := range admin.Providers {
		providerName := reflect.TypeOf(provider).String()

		// 包含字符串
		if find := strings.Contains(providerName, "*resources."); find {
			structName := strings.Replace(providerName, "*resources.", "", -1)

			if strings.ToLower(structName) == strings.ToLower(c.Params("resource")) {
				// 初始化实例
				resourceInstance = provider.(interface{ Init() interface{} }).Init()
			}
		}
	}

	return resourceInstance
}

// 模型
func (p *Quark) NewModel(resourceInstance interface{}) *gorm.DB {
	return resourceInstance.(interface{ NewModel(interface{}) *gorm.DB }).NewModel(resourceInstance)
}
