package requests

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin"
)

type Dashboard struct{}

// 资源
func (p *Dashboard) Resource(c *fiber.Ctx) interface{} {
	var dashboardInstance interface{}
	for _, provider := range admin.Providers {
		providerName := reflect.TypeOf(provider).String()

		// 包含字符串
		if find := strings.Contains(providerName, "*dashboards."); find {
			structName := strings.Replace(providerName, "*dashboards.", "", -1)
			if strings.ToLower(structName) == strings.ToLower(c.Params("dashboard")) {

				// 初始化实例
				dashboardInstance = provider.(interface{ Init() interface{} }).Init()
			}
		}
	}

	return dashboardInstance
}
