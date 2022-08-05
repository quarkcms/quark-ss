package admin

import (
	"reflect"
	"strings"

	"github.com/gobeam/stringy"
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/component/card"
	"github.com/quarkcms/quark-go/pkg/ui/component/form"
	"github.com/quarkcms/quark-go/pkg/ui/component/table"
	"github.com/quarkcms/quark-go/pkg/ui/component/tabs"
	"gorm.io/gorm"
)

// 结构体
type Resource struct {
	Layout
	Title        string
	SubTitle     string
	PerPage      interface{}
	IndexPolling int
	IndexOrder   string
	Model        interface{}
	Field        map[string]interface{}
	WithExport   bool
}

// 初始化
func (p *Resource) Init() *Resource {

	return p
}

// 获取模型
func (p *Resource) NewModel(resourceInstance interface{}) *gorm.DB {

	model := reflect.
		ValueOf(resourceInstance).
		Elem().
		FieldByName("Model").Interface()

	return (&db.Model{}).Model(&model)
}

// 列表页表格主体
func (p *Resource) IndexExtraRender(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	return nil
}

// 列表页工具栏
func (p *Resource) IndexToolBar(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	return (&table.ToolBar{}).Init().SetTitle(p.IndexTitle(c, resourceInstance)).SetActions(p.IndexActions(c, resourceInstance))
}

// 判断当前页面是否为列表页面 todo
func (p *Resource) IsIndex(c *fiber.Ctx) bool {

	uri := strings.Split(c.Path(), "/")

	return (uri[len(uri)-1] == "index")
}

//判断当前页面是否为创建页面
func (p *Resource) IsCreating(c *fiber.Ctx) bool {

	uri := strings.Split(c.Path(), "/")

	return (uri[len(uri)-1] == "create") || (uri[len(uri)-1] == "store")
}

// 判断当前页面是否为编辑页面
func (p *Resource) IsEditing(c *fiber.Ctx) bool {

	uri := strings.Split(c.Path(), "/")

	return (uri[len(uri)-1] == "edit") || (uri[len(uri)-1] == "update")
}

// 判断当前页面是否为详情页面
func (p *Resource) IsDetail(c *fiber.Ctx) bool {

	uri := strings.Split(c.Path(), "/")

	return (uri[len(uri)-1] == "detail")
}

// 判断当前页面是否为导出页面

func (p *Resource) isExport(c *fiber.Ctx) bool {

	uri := strings.Split(c.Path(), "/")

	return (uri[len(uri)-1] == "export")
}

// 列表标题
func (p *Resource) IndexTitle(c *fiber.Ctx, resourceInstance interface{}) string {
	return reflect.
		ValueOf(resourceInstance).
		Elem().
		FieldByName("Title").
		String() + "列表"
}

// 表单接口
func (p *Resource) FormApi(c *fiber.Ctx) string {
	return ""
}

// 创建表单的接口
func (p *Resource) CreationApi(c *fiber.Ctx, resourceInstance interface{}) string {

	formApi := resourceInstance.(interface{ FormApi(*fiber.Ctx) string }).FormApi(c)

	if formApi != "" {
		return formApi
	}

	uri := strings.Split(c.Path(), "/")

	if uri[len(uri)-1] == "index" {

		return stringy.New(stringy.New(c.Path()).ReplaceFirst("/api/", "")).ReplaceLast("/index", "/store")
	}

	return stringy.New(stringy.New(c.Path()).ReplaceFirst("/api/", "")).ReplaceLast("/create", "/store")
}

//更新表单的接口
func (p *Resource) UpdateApi(c *fiber.Ctx, resourceInstance interface{}) string {

	formApi := resourceInstance.(interface{ FormApi(*fiber.Ctx) string }).FormApi(c)

	if formApi != "" {
		return formApi
	}

	uri := strings.Split(c.Path(), "/")

	if uri[len(uri)-1] == "index" {

		return stringy.New(stringy.New(c.Path()).ReplaceFirst("/api/", "")).ReplaceLast("/index", "/save")
	}

	return stringy.New(stringy.New(c.Path()).ReplaceFirst("/api/", "")).ReplaceLast("/edit", "/save")
}

// 编辑页面获取表单数据接口
func (p *Resource) EditValueApi(c *fiber.Ctx) string {

	uri := strings.Split(c.Path(), "/")

	if uri[len(uri)-1] == "index" {

		return stringy.New(stringy.New(c.Path()).ReplaceFirst("/api/", "")).ReplaceLast("/index", "/edit/values?id=${id}")
	}

	return stringy.New(stringy.New(c.Path()).ReplaceFirst("/api/", "")).ReplaceLast("/edit", "/edit/values?id=${id}")
}

// 表单标题
func (p *Resource) FormTitle(c *fiber.Ctx, resourceInstance interface{}) string {
	value := reflect.ValueOf(resourceInstance).Elem()
	title := value.FieldByName("Title").String()

	if p.IsCreating(c) {
		return "创建" + title
	} else {
		if p.IsEditing(c) {
			return "编辑" + title
		}
	}

	return title
}

// 详情页标题

func (p *Resource) DetailTitle(c *fiber.Ctx, resourceInstance interface{}) string {
	value := reflect.ValueOf(resourceInstance).Elem()
	title := value.FieldByName("Title").String()

	return title + "详情"
}

// 列表页组件渲染
func (p *Resource) IndexComponentRender(c *fiber.Ctx, resourceInstance interface{}, data interface{}) interface{} {
	var component interface{}

	// 列表标题
	title := p.IndexTitle(c, resourceInstance)

	// 反射获取参数
	value := reflect.ValueOf(resourceInstance).Elem()
	indexPolling := value.FieldByName("IndexPolling").Int()

	// 列表页表格主体
	indexExtraRender := p.IndexExtraRender(c, resourceInstance)

	// 列表页工具栏
	indexToolBar := p.IndexToolBar(c, resourceInstance)

	// 列表页表格列
	indexColumns := p.IndexColumns(c, resourceInstance)

	// 列表页批量操作
	indexTableAlertActions := p.IndexTableAlertActions(c, resourceInstance)

	// 列表页搜索栏
	indexSearches := p.IndexSearches(c, resourceInstance)

	table := (&table.Component{}).
		Init().
		SetPolling(int(indexPolling)).
		SetTitle(title).
		SetTableExtraRender(indexExtraRender).
		SetToolBar(indexToolBar).
		SetColumns(indexColumns).
		SetBatchActions(indexTableAlertActions).
		SetSearches(indexSearches)

	// 获取分页
	perPage := reflect.
		ValueOf(resourceInstance).
		Elem().
		FieldByName("PerPage").Interface()

	// 不分页，直接返回数据
	if reflect.TypeOf(perPage).String() != "int" {
		component = table.SetDatasource(data)
	} else {
		current := data.(map[string]interface{})["currentPage"]
		perPage := data.(map[string]interface{})["perPage"]
		total := data.(map[string]interface{})["total"]
		items := data.(map[string]interface{})["items"]

		component = table.SetPagination(current.(int), perPage.(int), int(total.(int64)), 1).SetDatasource(items)
	}

	return component
}

// 渲染创建页组件
func (p *Resource) CreationComponentRender(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) interface{} {
	title := p.FormTitle(c, resourceInstance)
	formExtraActions := p.FormExtraActions(c, resourceInstance)
	api := p.CreationApi(c, resourceInstance)
	fields := p.CreationFieldsWithinComponents(c, resourceInstance)
	formActions := p.FormActions(c, resourceInstance)

	return p.FormComponentRender(
		c,
		resourceInstance,
		title,
		formExtraActions,
		api,
		fields,
		formActions,
		data,
	)
}

// 渲染编辑页组件
func (p *Resource) UpdateComponentRender(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) interface{} {
	title := p.FormTitle(c, resourceInstance)
	formExtraActions := p.FormExtraActions(c, resourceInstance)
	api := p.UpdateApi(c, resourceInstance)
	fields := p.UpdateFieldsWithinComponents(c, resourceInstance)
	formActions := p.FormActions(c, resourceInstance)

	return p.FormComponentRender(
		c,
		resourceInstance,
		title,
		formExtraActions,
		api,
		fields,
		formActions,
		data,
	)
}

// 渲染表单组件
func (p *Resource) FormComponentRender(
	c *fiber.Ctx,
	resourceInstance interface{},
	title string,
	extra interface{},
	api string,
	fields interface{},
	actions []interface{},
	data map[string]interface{}) interface{} {

	getFields, ok := fields.([]interface{})

	if ok {
		component := reflect.
			ValueOf(fields.([]interface{})[0]).
			Elem().
			FieldByName("Component").String()

		if component == "tabPane" {
			return p.FormWithinTabs(c, resourceInstance, title, extra, api, getFields, actions, data)
		} else {
			return p.FormWithinCard(c, resourceInstance, title, extra, api, fields, actions, data)
		}
	} else {
		return p.FormWithinCard(c, resourceInstance, title, extra, api, fields, actions, data)
	}
}

// 在卡片内的From组件
func (p *Resource) FormWithinCard(
	c *fiber.Ctx,
	resourceInstance interface{},
	title string,
	extra interface{},
	api string,
	fields interface{},
	actions []interface{},
	data map[string]interface{}) interface{} {

	formComponent := (&form.Component{}).
		Init().
		SetStyle(map[string]interface{}{
			"padding": "24px",
		}).
		SetApi(api).
		SetActions(actions).
		SetBody(fields).
		SetInitialValues(data)

	return (&card.Component{}).
		Init().
		SetTitle(title).
		SetHeaderBordered(true).
		SetExtra(extra).
		SetBody(formComponent)
}

// 在标签页内的From组件
func (p *Resource) FormWithinTabs(
	c *fiber.Ctx,
	resourceInstance interface{},
	title string,
	extra interface{},
	api string,
	fields interface{},
	actions []interface{},
	data map[string]interface{}) interface{} {

	tabsComponent := (&tabs.Component{}).Init().SetTabPanes(fields).SetTabBarExtraContent(extra)

	return (&form.Component{}).
		Init().
		SetStyle(map[string]interface{}{
			"backgroundColor": "#fff",
			"paddingBottom":   "20px",
		}).
		SetApi(api).
		SetActions(actions).
		SetBody(tabsComponent).
		SetInitialValues(data)
}

// 渲染详情页组件
func (p *Resource) DetailComponentRender(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) interface{} {
	title := p.DetailTitle(c, resourceInstance)
	formExtraActions := p.DetailExtraActions(c, resourceInstance)
	fields := p.DetailFieldsWithinComponents(c, resourceInstance, data)
	formActions := p.DetailActions(c, resourceInstance)

	return p.DetailWithinCard(
		c,
		resourceInstance,
		title,
		formExtraActions,
		fields,
		formActions,
		data,
	)
}

// 在卡片内的详情页组件
func (p *Resource) DetailWithinCard(
	c *fiber.Ctx,
	resourceInstance interface{},
	title string,
	extra interface{},
	fields interface{},
	actions []interface{},
	data map[string]interface{}) interface{} {

	return (&card.Component{}).
		Init().
		SetTitle(title).
		SetHeaderBordered(true).
		SetExtra(extra).
		SetBody(fields)
}

// 在标签页内的详情页组件
func (p *Resource) DetailWithinTabs(
	c *fiber.Ctx,
	resourceInstance interface{},
	title string,
	extra interface{},
	fields interface{},
	actions []interface{},
	data map[string]interface{}) interface{} {

	return (&tabs.Component{}).Init().SetTabPanes(fields).SetTabBarExtraContent(extra)
}

// 设置单列字段
func (p *Resource) SetField(fieldData map[string]interface{}) interface{} {
	p.Field = fieldData

	return p
}

// 列表页面显示前回调
func (p *Resource) BeforeIndexShowing(c *fiber.Ctx, list []map[string]interface{}) []interface{} {
	result := []interface{}{}

	for _, v := range list {
		result = append(result, v)
	}

	return result
}

// 详情页页面显示前回调
func (p *Resource) BeforeDetailShowing(c *fiber.Ctx, data map[string]interface{}) map[string]interface{} {
	return data
}

// 创建页面显示前回调
func (p *Resource) BeforeCreating(c *fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{}
}

// 编辑页面显示前回调
func (p *Resource) BeforeEditing(c *fiber.Ctx, data map[string]interface{}) map[string]interface{} {
	return data
}

// 保存数据前回调
func (p *Resource) BeforeSaving(c *fiber.Ctx, submitData map[string]interface{}) interface{} {
	return submitData
}

// 保存数据后回调
func (p *Resource) AfterSaved(c *fiber.Ctx, model *gorm.DB) interface{} {
	return model
}

// 数据导出前回调
func (p *Resource) BeforeExporting(c *fiber.Ctx, list []map[string]interface{}) []interface{} {
	result := []interface{}{}

	for _, v := range list {
		result = append(result, v)
	}

	return result
}

// 数据导入前回调
func (p *Resource) BeforeImporting(c *fiber.Ctx, list [][]interface{}) [][]interface{} {
	return list
}
