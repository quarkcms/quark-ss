package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/controllers"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/middleware"
)

type Admin struct{}

// 路由
func (p *Admin) Route(app *fiber.App) {
	ag := app.Group("/api/admin")
	ag.Get("/login", (&controllers.Login{}).Show)
	ag.Post("/login", (&controllers.Login{}).Login)
	ag.Get("/logout", (&controllers.Login{}).Logout)
	ag.Get("/captcha", (&controllers.Captcha{}).Make)

	amg := app.Group("/api/admin", (&middleware.AdminMiddleware{}).Handle)
	amg.Get("/dashboard/:dashboard", (&controllers.Dashboard{}).Handle)             // 仪表盘
	amg.Get("/:resource/index", (&controllers.ResourceIndex{}).Handle)              // 列表页面
	amg.Get("/:resource/editable", (&controllers.ResourceEditable{}).Handle)        // 表格行内编辑
	amg.All("/:resource/action/:uriKey", (&controllers.ResourceAction{}).Handle)    // 执行行为
	amg.Get("/:resource/create", (&controllers.ResourceCreate{}).Handle)            // 创建页面
	amg.Post("/:resource/store", (&controllers.ResourceStore{}).Handle)             // 创建方法
	amg.Get("/:resource/edit", (&controllers.ResourceEdit{}).Handle)                // 编辑页面
	amg.Get("/:resource/edit/values", (&controllers.ResourceEdit{}).Values)         // 获取编辑页面表单值
	amg.Post("/:resource/save", (&controllers.ResourceUpdate{}).Handle)             // 保存编辑值
	amg.Get("/:resource/detail", (&controllers.ResourceDetail{}).Handle)            // 详情页面
	amg.Get("/:resource/export", (&controllers.ResourceExport{}).Handle)            // 导出
	amg.All("/:resource/import", (&controllers.ResourceImport{}).Handle)            // 导入
	amg.Get("/:resource/import/template", (&controllers.ResourceImport{}).Template) // 导入模板

	// 通用表单资源
	amg.Get("/:resource/:uriKey-form", (&controllers.ResourceCreate{}).Handle)
	amg.Get("/:resource/:uriKey/form", (&controllers.ResourceCreate{}).Handle)

	// 直接加载前端组件
	amg.Get("/pages/:component", (&controllers.Pages{}).Handle)

	// 图片上传、下载
	amg.Get("/picture/getLists", (&controllers.Picture{}).GetLists)
	amg.Post("/picture/upload", (&controllers.Picture{}).Upload)
	amg.Get("/picture/download", (&controllers.Picture{}).Download)
	amg.All("/picture/delete", (&controllers.Picture{}).Delete)
	amg.Post("/picture/crop", (&controllers.Picture{}).Crop)

	// 文件上传、下载
	amg.Post("/file/upload", (&controllers.File{}).Upload)
	amg.Get("/file/download", (&controllers.File{}).Download)
}
