package admin

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/component/descriptions"
	"github.com/quarkcms/quark-go/pkg/ui/component/table"
	"github.com/quarkcms/quark-go/pkg/ui/component/tabs"
)

// 列表字段
func (p *Resource) IndexFields(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	fields := p.getFields(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface {
			IsShownOnIndex() bool
		}); ok {
			isShownOnIndex := v.IsShownOnIndex()

			if isShownOnIndex {
				items = append(items, v)
			}
		}
	}

	return items
}

// 表格列
func (p *Resource) IndexColumns(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	fields := p.IndexFields(c, resourceInstance)
	var columns []interface{}

	for _, v := range fields.([]interface{}) {
		isShownOnIndex := v.(interface {
			IsShownOnIndex() bool
		}).IsShownOnIndex()

		if isShownOnIndex {
			getColumn := p.fieldToColumn(c, v)

			if getColumn != nil {
				columns = append(columns, getColumn)
			}
		}
	}

	// 行内行为
	indexTableRowActions := resourceInstance.(interface {
		IndexTableRowActions(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).IndexTableRowActions(c, resourceInstance)

	if len(indexTableRowActions.([]interface{})) > 0 {
		column := (&table.Column{}).
			Init().
			SetTitle("操作").
			SetAttribute("action").
			SetValueType("option").
			SetActions(indexTableRowActions).
			SetFixed("right")

		columns = append(columns, column)
	}

	return columns
}

// 将表单项转换为表格列
func (p *Resource) fieldToColumn(c *fiber.Ctx, field interface{}) interface{} {

	reflectElem := reflect.
		ValueOf(field).
		Elem()

	// 字段
	name := reflectElem.
		FieldByName("Name").
		String()

	// 文字
	label := reflectElem.
		FieldByName("Label").
		String()

	// 组件类型
	component := reflectElem.
		FieldByName("Component").
		String()

	// 是否可编辑
	editable := reflectElem.
		FieldByName("Editable").
		Bool()

	// 是否可编辑
	getColumn := reflectElem.
		FieldByName("Column").
		Interface()

	column := getColumn.(*table.Column).
		SetTitle(label).
		SetAttribute(name)

	switch component {
	case "idField":
		// 是否显示在列表
		onIndexDisplayed := reflectElem.
			FieldByName("OnIndexDisplayed").
			Bool()

		if onIndexDisplayed {
			column = column.SetValueType("text")
		} else {
			return nil
		}
	case "hiddenField":
		return nil
	case "textField":
		column = column.SetValueType("text")
	case "textAreaField":
		column = column.SetValueType("text")
	case "selectField":
		valueEnum := field.(interface {
			GetValueEnum() map[interface{}]interface{}
		}).GetValueEnum()
		column = column.SetValueType("select").SetValueEnum(valueEnum)
	case "radioField":
		valueEnum := field.(interface {
			GetValueEnum() map[interface{}]interface{}
		}).GetValueEnum()
		column = column.SetValueType("radio").SetValueEnum(valueEnum)
	case "switchField":
		valueEnum := field.(interface {
			GetSwitchValueEnum() map[interface{}]interface{}
		}).GetSwitchValueEnum()
		column = column.SetValueType("select").SetValueEnum(valueEnum)
	case "imageField":
		column = column.SetValueType("image")
	default:
		column = column.SetValueType(component)
	}

	if editable {
		// 可编辑，设置编辑
		options := reflectElem.
			FieldByName("Options").
			Interface()

		// 可编辑api地址
		editableApi := strings.Replace(strings.Replace(c.Path(), "/api/", "", -1), "/index", "/editable", -1)

		// 设置编辑项
		column = column.SetEditable(component, options, editableApi)
	}

	return column
}

// 创建页字段
func (p *Resource) CreationFields(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := p.getFields(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnCreation := v.(interface {
			IsShownOnCreation() bool
		}).IsShownOnCreation()

		if isShownOnCreation {
			items = append(items, v)
		}
	}

	return items

}

// 不包含When组件内字段的创建页字段
func (p *Resource) CreationFieldsWithoutWhen(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := p.getFieldsWithoutWhen(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnCreation := v.(interface {
			IsShownOnCreation() bool
		}).IsShownOnCreation()

		if isShownOnCreation {
			items = append(items, v)
		}
	}

	return items
}

// 包裹在组件内的创建页字段
func (p *Resource) CreationFieldsWithinComponents(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := resourceInstance.(interface {
		Fields(c *fiber.Ctx) []interface{}
	}).Fields(c)
	var items []interface{}

	for _, v := range fields {

		component := reflect.
			ValueOf(v).
			Elem().
			FieldByName("Component").String()

		if component == "tabPane" {

			body := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").Interface()

			var subItems []interface{}
			for _, sv := range body.([]interface{}) {
				isShownOnCreation := sv.(interface {
					IsShownOnCreation() bool
				}).IsShownOnCreation()

				if isShownOnCreation {
					sv.(interface {
						BuildFrontendRules(c *fiber.Ctx) interface{}
					}).BuildFrontendRules(c)

					subItems = append(subItems, sv)
				}
			}

			v.(interface {
				SetBody(interface{}) *tabs.TabPane
			}).SetBody(subItems)

			items = append(items, v)
		} else {
			isShownOnCreation := v.(interface {
				IsShownOnCreation() bool
			}).IsShownOnCreation()

			if isShownOnCreation {
				v.(interface {
					BuildFrontendRules(c *fiber.Ctx) interface{}
				}).BuildFrontendRules(c)
				items = append(items, v)
			}
		}
	}

	return items
}

// 编辑页字段
func (p *Resource) UpdateFields(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	fields := p.getFields(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnUpdate := v.(interface {
			IsShownOnUpdate() bool
		}).IsShownOnUpdate()

		if isShownOnUpdate {
			items = append(items, v)
		}
	}

	return items
}

// 不包含When组件内字段的编辑页字段
func (p *Resource) UpdateFieldsWithoutWhen(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := p.getFieldsWithoutWhen(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnUpdate := v.(interface {
			IsShownOnUpdate() bool
		}).IsShownOnUpdate()

		if isShownOnUpdate {
			items = append(items, v)
		}
	}

	return items
}

// 包裹在组件内的编辑页字段
func (p *Resource) UpdateFieldsWithinComponents(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := resourceInstance.(interface {
		Fields(c *fiber.Ctx) []interface{}
	}).Fields(c)
	var items []interface{}

	for _, v := range fields {

		component := reflect.
			ValueOf(v).
			Elem().
			FieldByName("Component").String()

		if component == "tabPane" {

			body := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").Interface()

			var subItems []interface{}
			for _, sv := range body.([]interface{}) {
				isShownOnUpdate := sv.(interface {
					IsShownOnUpdate() bool
				}).IsShownOnUpdate()

				if isShownOnUpdate {

					sv.(interface {
						BuildFrontendRules(c *fiber.Ctx) interface{}
					}).BuildFrontendRules(c)

					subItems = append(subItems, sv)
				}
			}

			v.(interface {
				SetBody(interface{}) *tabs.TabPane
			}).SetBody(subItems)

			items = append(items, v)
		} else {
			isShownOnUpdate := v.(interface {
				IsShownOnUpdate() bool
			}).IsShownOnUpdate()

			if isShownOnUpdate {

				v.(interface {
					BuildFrontendRules(c *fiber.Ctx) interface{}
				}).BuildFrontendRules(c)

				items = append(items, v)
			}
		}
	}

	return items
}

// 详情页字段
func (p *Resource) DetailFields(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	fields := p.getFields(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnDetail := v.(interface {
			IsShownOnDetail() bool
		}).IsShownOnDetail()

		if isShownOnDetail {
			items = append(items, v)
		}
	}

	return items
}

// 包裹在组件内的详情页字段
func (p *Resource) DetailFieldsWithinComponents(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) interface{} {
	componentType := "description"

	fields := resourceInstance.(interface {
		Fields(c *fiber.Ctx) []interface{}
	}).Fields(c)
	var items []interface{}

	for _, v := range fields {

		component := reflect.
			ValueOf(v).
			Elem().
			FieldByName("Component").String()

		if component == "tabPane" {

			body := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").Interface()

			var subItems []interface{}
			for _, sv := range body.([]interface{}) {
				isShownOnDetail := sv.(interface {
					IsShownOnDetail() bool
				}).IsShownOnDetail()

				if isShownOnDetail {
					getColumn := p.fieldToColumn(c, sv)
					subItems = append(subItems, getColumn)
				}
			}

			descriptions := (&descriptions.Component{}).Init().SetStyle(map[string]interface{}{
				"padding": "24px",
			}).
				SetTitle("").
				SetColumn(2).
				SetColumns(subItems).
				SetDataSource(data).
				SetActions(p.DetailActions(c, resourceInstance))

			v.(interface{ SetBody(interface{}) interface{} }).SetBody(descriptions)
			items = append(items, v)
		} else {
			isShownOnDetail := v.(interface {
				IsShownOnDetail() bool
			}).IsShownOnDetail()

			if isShownOnDetail {
				getColumn := p.fieldToColumn(c, v)
				items = append(items, getColumn)
			}
		}
	}

	if componentType == "description" {
		return (&descriptions.Component{}).
			Init().
			SetStyle(map[string]interface{}{
				"padding": "24px",
			}).
			SetTitle("").
			SetColumn(2).
			SetColumns(items).
			SetDataSource(data).
			SetActions(p.DetailActions(c, resourceInstance))
	} else {
		return items
	}
}

// 导出字段
func (p *Resource) ExportFields(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := p.getFields(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnExport := v.(interface {
			IsShownOnExport() bool
		}).IsShownOnExport()

		if isShownOnExport {
			items = append(items, v)
		}
	}

	return items
}

// 导入字段
func (p *Resource) ImportFields(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := p.getFields(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnImport := v.(interface {
			IsShownOnImport() bool
		}).IsShownOnImport()

		if isShownOnImport {
			items = append(items, v)
		}
	}

	return items
}

// 不包含When组件内字段的导入字段
func (p *Resource) ImportFieldsWithoutWhen(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := p.getFieldsWithoutWhen(c, resourceInstance)
	var items []interface{}

	for _, v := range fields.([]interface{}) {

		isShownOnImport := v.(interface {
			IsShownOnImport() bool
		}).IsShownOnImport()

		if isShownOnImport {
			items = append(items, v)
		}
	}

	return items
}

// 获取字段
func (p *Resource) getFields(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := resourceInstance.(interface {
		Fields(c *fiber.Ctx) []interface{}
	}).Fields(c)

	return p.findFields(fields, true)
}

// 获取不包含When组件的字段
func (p *Resource) getFieldsWithoutWhen(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	fields := resourceInstance.(interface {
		Fields(c *fiber.Ctx) []interface{}
	}).Fields(c)

	return p.findFields(fields, false)
}

// 查找字段
func (p *Resource) findFields(fields interface{}, when bool) interface{} {
	var items []interface{}

	if fields, ok := fields.([]interface{}); ok {
		for _, v := range fields {
			hasBody := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").IsValid()

			if hasBody {
				body := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").Interface()

				getItems := p.findFields(body, true)

				if getItems, ok := getItems.([]interface{}); ok {
					if len(getItems) > 0 {
						items = append(items, getItems...)
					}
				}

			} else {

				component := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Component").String()

				if strings.Contains(component, "Field") {
					items = append(items, v)
					if when {
						whenFields := p.getWhenFields(v)
						if len(whenFields) > 0 {
							items = append(items, whenFields...)
						}
					}
				}
			}
		}
	}

	return items
}

// 获取When组件中的字段
func (p *Resource) getWhenFields(item interface{}) []interface{} {

	var items []interface{}

	when := reflect.
		ValueOf(item).
		Elem().
		FieldByName("When").Interface()

	if when == nil {
		return items
	}

	whenItems := reflect.
		ValueOf(when).
		Elem().
		FieldByName("Items").Interface()

	if whenItems == nil {
		return items
	}

	whenItems, ok := whenItems.([]interface{})

	if ok {
		for _, v := range whenItems.([]interface{}) {
			body := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").Interface()

			if body != nil {

				if body, ok := body.([]interface{}); ok {
					if len(body) > 0 {
						items = append(items, body...)
					}
				}

				if body, ok := body.(interface{}); ok {
					items = append(items, body)
				}
			}
		}
	}

	return items
}
