package admin

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/config"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"github.com/quarkcms/quark-go/pkg/ui/component/footer"
	"github.com/quarkcms/quark-go/pkg/ui/component/layout"
	"github.com/quarkcms/quark-go/pkg/ui/component/page"
	"github.com/quarkcms/quark-go/pkg/ui/component/pagecontainer"
)

// 结构体
type Layout struct{}

// 页面组件渲染
func (p *Layout) PageComponentRender(c *fiber.Ctx, resourceInstance interface{}, content interface{}) interface{} {
	layoutComponent := resourceInstance.(interface {
		LayoutComponentRender(c *fiber.Ctx, resourceInstance interface{}, content interface{}) interface{}
	}).LayoutComponentRender(c, resourceInstance, content)

	return (&page.Component{}).
		Init().
		SetStyle(map[string]interface{}{
			"height": "100vh",
		}).
		SetBody(layoutComponent).
		JsonSerialize()
}

// 页面布局组件渲染
func (p *Layout) LayoutComponentRender(c *fiber.Ctx, resourceInstance interface{}, content interface{}) interface{} {

	// 获取登录管理员信息
	adminId := utils.Admin(c, "id")

	// 获取管理员菜单
	getMenus := (&models.Admin{}).GetMenus(adminId.(int))

	// 页脚
	footer := (&footer.Component{}).
		Init().
		SetCopyright(config.Get("admin.copyright").(string)).
		SetLinks(config.Get("admin.links").([]map[string]interface{}))

	return (&layout.Component{}).
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
		SetBody(resourceInstance.(interface {
			PageContainerComponentRender(c *fiber.Ctx, resourceInstance interface{}, content interface{}) interface{}
		}).PageContainerComponentRender(c, resourceInstance, content)).
		SetFooter(footer)
}

// 页面容器组件渲染
func (p *Layout) PageContainerComponentRender(c *fiber.Ctx, resourceInstance interface{}, content interface{}) interface{} {

	value := reflect.ValueOf(resourceInstance).Elem()
	title := value.FieldByName("Title").String()
	subTitle := value.FieldByName("SubTitle").String()

	header := (&pagecontainer.PageHeader{}).
		Init().
		SetTitle(title).
		SetSubTitle(subTitle)

	return (&pagecontainer.Component{}).Init().SetHeader(header).SetBody(content)
}

// 组件渲染
func (p *Layout) Render(c *fiber.Ctx, resourceInstance interface{}, content interface{}) interface{} {

	return resourceInstance.(interface {
		PageComponentRender(c *fiber.Ctx, resourceInstance interface{}, content interface{}) interface{}
	}).PageComponentRender(c, resourceInstance, content)
}
