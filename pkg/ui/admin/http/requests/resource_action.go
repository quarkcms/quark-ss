package requests

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"gorm.io/gorm"
)

type ResourceAction struct {
	Quark
}

// 执行行为
func (p *ResourceAction) HandleAction(c *fiber.Ctx) error {
	var result error
	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)

	id := c.Query("id")

	if id != "" {
		if strings.Contains(id, ",") {
			model.Where("id IN ?", strings.Split(id, ","))
		} else {
			model.Where("id = ?", id)
		}
	}

	// hack:自动生成权限信息
	if len(utils.GetPermissions()) <= 0 {
		utils.SetPermissions(resourceToPermission(c))
	}

	actions := resourceInstance.(interface {
		Actions(c *fiber.Ctx) []interface{}
	}).Actions(c)

	for _, v := range actions {

		// uri唯一标识
		uriKey := v.(interface {
			GetUriKey(interface{}) string
		}).GetUriKey(v)

		actionType := v.(interface{ GetActionType() string }).GetActionType()

		if actionType == "dropdown" {
			dropdownActions := v.(interface{ GetActions() []interface{} }).GetActions()
			for _, dropdownAction := range dropdownActions {
				// uri唯一标识
				uriKey := dropdownAction.(interface {
					GetUriKey(interface{}) string
				}).GetUriKey(dropdownAction)

				if c.Params("uriKey") == uriKey {
					result = dropdownAction.(interface {
						Handle(*fiber.Ctx, *gorm.DB) error
					}).Handle(c, model)
				}
			}
		} else {
			if c.Params("uriKey") == uriKey {
				result = v.(interface {
					Handle(*fiber.Ctx, *gorm.DB) error
				}).Handle(c, model)
			}
		}
	}

	return result
}

// 资源转换成权限
func resourceToPermission(c *fiber.Ctx) []string {
	permissions := []string{}
	routes := []string{
		"api/admin/dashboard/:dashboard",
		"api/admin/:resource/index",
		"api/admin/:resource/editable",
		"api/admin/:resource/action/:uriKey",
		"api/admin/:resource/create",
		"api/admin/:resource/store",
		"api/admin/:resource/edit",
		"api/admin/:resource/edit/values",
		"api/admin/:resource/save",
		"api/admin/:resource/detail",
	}

	for _, provider := range admin.Providers {
		providerName := reflect.TypeOf(provider).String()

		// 处理仪表盘
		if find := strings.Contains(providerName, "*dashboards."); find {
			structName := strings.Replace(providerName, "*dashboards.", "", -1)
			for _, v := range routes {
				if strings.Contains(v, ":dashboard") {
					permissions = append(permissions, strings.Replace(v, ":dashboard", strings.ToLower(structName), -1))
				}
			}
		}

		// 处理资源
		if find := strings.Contains(providerName, "*resources."); find {
			structName := strings.Replace(providerName, "*resources.", "", -1)
			for _, v := range routes {
				if strings.Contains(v, ":resource") {
					v = strings.Replace(v, ":resource", strings.ToLower(structName), -1)

					//处理行为
					if strings.Contains(v, ":uriKey") {

						// 初始化实例
						resourceInstance := provider.(interface{ Init() interface{} }).Init()
						actions := resourceInstance.(interface {
							Actions(c *fiber.Ctx) []interface{}
						}).Actions(c)

						for _, av := range actions {

							// uri唯一标识
							uriKey := av.(interface {
								GetUriKey(interface{}) string
							}).GetUriKey(av)

							actionType := av.(interface{ GetActionType() string }).GetActionType()

							if actionType == "dropdown" {
								dropdownActions := av.(interface{ GetActions() []interface{} }).GetActions()
								for _, dropdownAction := range dropdownActions {
									// uri唯一标识
									uriKey := dropdownAction.(interface {
										GetUriKey(interface{}) string
									}).GetUriKey(dropdownAction)

									v = strings.Replace(v, ":uriKey", uriKey, -1)
								}
							} else {
								v = strings.Replace(v, ":uriKey", uriKey, -1)
							}
						}
					}

					permissions = append(permissions, v)
				}
			}
		}
	}

	return permissions
}
