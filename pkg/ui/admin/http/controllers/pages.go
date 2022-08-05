package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/config"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"github.com/quarkcms/quark-go/pkg/ui/component/footer"
	"github.com/quarkcms/quark-go/pkg/ui/component/layout"
	"github.com/quarkcms/quark-go/pkg/ui/component/page"
)

type Pages struct{}

// 执行行为
func (p *Pages) Handle(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"component": c.Params("component"),
	}

	// 获取登录管理员信息
	adminId := utils.Admin(c, "id")

	// 获取管理员菜单
	getMenus := (&models.Admin{}).GetMenus(adminId.(int))

	// 页脚
	footer := (&footer.Component{}).
		Init().
		SetCopyright(config.Get("admin.copyright").(string)).
		SetLinks(config.Get("admin.links").([]map[string]interface{}))

	layoutComponent := (&layout.Component{}).
		Init().
		SetTitle(config.Get("admin.name").(string)).
		SetLogo(config.Get("admin.layout.logo")).
		SetHeaderActions(config.Get("admin.layout.header_actions").([]map[string]interface{})).
		SetLayout(config.Get("admin.layout.layout").(string)).
		SetSplitMenus(config.Get("admin.layout.split_menus").(bool)).
		SetHeaderTheme(config.Get("admin.layout.header_theme").(string)).
		SetContentWidth(config.Get("admin.layout.content_width").(string)).
		SetNavTheme(config.Get("admin.layout.nav_theme").(string)).
		SetPrimaryColor(config.Get("admin.layout.primary_color").(string)).
		SetFixSiderbar(config.Get("admin.layout.fix_siderbar").(bool)).
		SetFixedHeader(config.Get("admin.layout.fixed_header").(bool)).
		SetIconfontUrl(config.Get("admin.layout.iconfont_url").(string)).
		SetLocale(config.Get("admin.layout.locale").(string)).
		SetSiderWidth(config.Get("admin.layout.sider_width").(int)).
		SetMenu(getMenus).
		SetBody(data).
		SetFooter(footer)

	pageComponent := (&page.Component{}).
		Init().
		SetStyle(map[string]interface{}{
			"height": "100vh",
		}).
		SetBody(layoutComponent).
		JsonSerialize()

	return c.JSON(pageComponent)
}
