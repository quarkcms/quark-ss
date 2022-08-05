package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/framework/rand"
	"github.com/quarkcms/quark-go/pkg/ui/admin/http/requests"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"github.com/xuri/excelize/v2"
)

type ResourceExport struct{}

// 导出数据
func (p *ResourceExport) Handle(c *fiber.Ctx) error {

	resourceExport := &requests.ResourceExport{}

	// 资源实例
	resourceInstance := resourceExport.Resource(c)

	if resourceInstance == nil {
		return c.SendStatus(404)
	}

	data := resourceExport.IndexQuery(c)

	// 获取列表字段
	fields := resourceInstance.(interface {
		ExportFields(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).ExportFields(c, resourceInstance)

	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")

	rowData := map[string]interface{}{}

	var a = 'a'
	for _, fieldValue := range fields.([]interface{}) {
		Label := reflect.
			ValueOf(fieldValue).
			Elem().
			FieldByName("Label").
			String()

		f.SetCellValue("Sheet1", string(a)+"1", Label)

		a++
	}

	for dataKey, dataValue := range data.([]interface{}) {
		var a = 'a'
		for _, fieldValue := range fields.([]interface{}) {

			name := reflect.
				ValueOf(fieldValue).
				Elem().
				FieldByName("Name").String()

			component := reflect.
				ValueOf(fieldValue).
				Elem().
				FieldByName("Component").String()

			switch component {
			case "inputNumberField":
				rowData[name] = dataValue.(map[string]interface{})[name]

			case "textField":
				rowData[name] = dataValue.(map[string]interface{})[name]

			case "selectField":
				options := reflect.
					ValueOf(fieldValue).
					Elem().
					FieldByName("Options").Interface()

				rowData[name] = p.getOptionValue(options, dataValue.(map[string]interface{})[name])

			case "cascaderField":
				options := reflect.
					ValueOf(fieldValue).
					Elem().
					FieldByName("Options").Interface()

				rowData[name] = p.getOptionValue(options, dataValue.(map[string]interface{})[name])

			case "checkboxField":
				options := reflect.
					ValueOf(fieldValue).
					Elem().
					FieldByName("Options").Interface()

				rowData[name] = p.getOptionValue(options, dataValue.(map[string]interface{})[name])

			case "radioField":
				options := reflect.
					ValueOf(fieldValue).
					Elem().
					FieldByName("Options").Interface()

				rowData[name] = p.getOptionValue(options, dataValue.(map[string]interface{})[name])

			case "switchField":
				options := reflect.
					ValueOf(fieldValue).
					Elem().
					FieldByName("Options").Interface()

				rowData[name] = p.getSwitchValue(options, dataValue.(map[string]interface{})[name].(int))

			default:
				rowData[name] = dataValue.(map[string]interface{})[name]
			}

			f.SetCellValue("Sheet1", string(a)+strconv.Itoa(dataKey+2), rowData[name])
			a++
		}
	}

	f.SetActiveSheet(index)

	filePath := "./storage/app/public/exports/"
	fileName := rand.MakeAlphanumeric(40) + ".xlsx"

	// 不存在路径，则创建
	if utils.PathExist(filePath) == false {
		err := os.MkdirAll(filePath, 0666)
		if err != nil {
			return msg.Error(err.Error(), "")
		}
	}

	if err := f.SaveAs(filePath + fileName); err != nil {
		fmt.Println(err)
	}

	return c.Redirect(c.BaseURL() + strings.Replace(filePath+fileName, "./storage/app/public", "/storage", -1))
}

// 获取属性值
func (p *ResourceExport) getOptionValue(options interface{}, value interface{}) string {
	result := ""
	arr := []interface{}{}

	if value, ok := value.(string); ok {
		if strings.Contains(value, "[") || strings.Contains(value, "{") {
			json.Unmarshal([]byte(value), &arr)
		}
	}

	if len(arr) > 0 {
		for _, option := range options.([]interface{}) {
			for _, v := range arr {
				if v == option.(map[string]interface{})["value"] {
					result = result + option.(map[string]interface{})["label"].(string)
				}
			}
		}
	} else {

		for _, option := range options.([]interface{}) {
			if value.(string) == option.(map[string]interface{})["value"] {
				result = option.(map[string]interface{})["label"].(string)
			}
		}
	}

	return result
}

// 获取开关组件值
func (p *ResourceExport) getSwitchValue(options interface{}, value int) string {
	if value == 1 {
		return options.(map[string]interface{})["on"].(string)
	} else {
		return options.(map[string]interface{})["off"].(string)
	}
}
