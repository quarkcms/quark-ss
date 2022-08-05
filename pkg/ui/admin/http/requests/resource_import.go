package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/rand"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"github.com/quarkcms/quark-go/pkg/ui/component/space"
	"github.com/quarkcms/quark-go/pkg/ui/component/tpl"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ResourceImport struct {
	Quark
}

// 文件结构体
type File struct {
	Id             int
	ObjType        string
	ObjId          int
	FileCategoryId int
	Sort           int
	Name           string
	Size           string
	Ext            string
	Path           string
	Md5            string
	Status         bool
}

// 执行行为
func (p *ResourceImport) HandleImport(c *fiber.Ctx) (interface{}, bool, error) {
	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)
	fileId := data["fileId"]
	getFileId := 0

	if value, ok := fileId.(map[string]interface{}); ok {
		if value, ok := value["0"]; ok {
			if value, ok := value.(map[string]interface{})["id"]; ok {
				getFileId = int(value.(float64))
			}
		}
	}

	if getFileId == 0 {
		return nil, false, errors.New("参数错误！")
	}

	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)

	importData, err := utils.Import(getFileId)

	if err != nil {
		return nil, false, err
	}

	// 表格头部
	importHead := importData[0]

	// 去除表格头部
	importData = importData[1:]

	// 导入前回调
	lists := resourceInstance.(interface {
		BeforeImporting(c *fiber.Ctx, list [][]interface{}) [][]interface{}
	}).BeforeImporting(c, importData)

	importResult := true
	importTotalNum := len(lists)
	importSuccessedNum := 0
	importFailedNum := 0
	importFailedData := [][]interface{}{}

	// 获取字段
	fields := resourceInstance.(interface {
		ImportFields(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).ImportFields(c, resourceInstance)

	for _, item := range lists {

		formValues := p.transformFormValues(fields, item)

		validator := resourceInstance.(interface {
			ValidatorForImport(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) error
		}).ValidatorForImport(c, resourceInstance, formValues)

		if validator != nil {

			importResult = false
			importFailedNum = importFailedNum + 1  // 记录错误条数
			item = append(item, validator.Error()) // 记录错误信息

			//记录错误数据
			importFailedData = append(importFailedData, item)
		}

		submitData := resourceInstance.(interface {
			BeforeSaving(c *fiber.Ctx, data map[string]interface{}) interface{}
		}).BeforeSaving(c, formValues) // 保存前回调

		if value, ok := submitData.(error); ok {
			importResult = false
			importFailedNum = importFailedNum + 1 // 记录错误条数
			item = append(item, value.Error())    // 记录错误信息

			//记录错误数据
			importFailedData = append(importFailedData, item)
		}

		data := p.getSubmitData(
			fields,
			submitData,
		)

		// 获取对象
		getModel := model.Create(data)

		if getModel.Error != nil {

			importResult = false
			importFailedNum = importFailedNum + 1       // 记录错误条数
			item = append(item, getModel.Error.Error()) // 记录错误信息

			//记录错误数据
			importFailedData = append(importFailedData, item)
		} else {
			resourceInstance.(interface {
				AfterSaved(c *fiber.Ctx, model *gorm.DB) interface{}
			}).AfterSaved(c, getModel)

			importSuccessedNum = importSuccessedNum + 1
		}
	}

	if importResult {
		return nil, true, nil
	} else {
		importHead = append(importHead, "错误信息")

		f := excelize.NewFile()
		// Create a new sheet.
		index := f.NewSheet("Sheet1")

		//定义一个字符 变量a 是一个byte类型的 表示单个字符
		var a = 'a'

		//生成26个字符
		for i := 1; i <= len(importHead); i++ {
			// 设置单元格的值
			f.SetCellValue("Sheet1", string(a)+"1", importHead[i-1])
			a++
		}

		for k, v := range importFailedData {
			//定义一个字符 变量a 是一个byte类型的 表示单个字符
			var a = 'a'

			//生成26个字符
			for i := 1; i <= len(v); i++ {
				// 设置单元格的值
				f.SetCellValue("Sheet1", string(a)+strconv.Itoa(k+2), v[i-1])
				a++
			}
		}

		f.SetActiveSheet(index)

		filePath := "./storage/app/public/failImports/"
		fileName := rand.MakeAlphanumeric(40) + ".xlsx"

		// 不存在路径，则创建
		if utils.PathExist(filePath) == false {
			err := os.MkdirAll(filePath, 0666)
			if err != nil {
				return nil, false, err
			}
		}

		if err := f.SaveAs(filePath + fileName); err != nil {
			fmt.Println(err)
		}

		importTotalNumTpl := (&tpl.Component{}).
			Init().
			SetBody("导入总量: " + strconv.Itoa(importTotalNum))

		importSuccessedNumTpl := (&tpl.Component{}).
			Init().
			SetBody("成功数量: " + strconv.Itoa(importSuccessedNum))

		importFailedNumTpl := (&tpl.Component{}).
			Init().
			SetBody("失败数量: <span style='color:#ff4d4f'>" + strconv.Itoa(importFailedNum) + "</span> <a href='" + c.BaseURL() + "/storage/failImports/" + fileName + "' target='_blank'>下载失败数据</a>")

		component := (&space.Component{}).
			Init().
			SetBody([]interface{}{
				importTotalNumTpl,
				importSuccessedNumTpl,
				importFailedNumTpl,
			}).
			SetDirection("vertical").
			SetSize("small").
			SetStyle(map[string]interface{}{
				"marginLeft":   "50px",
				"marginBottom": "20px",
			})

		return component, false, nil
	}
}

// 将表格数据转换成表单数据
func (p *ResourceImport) transformFormValues(fields interface{}, data []interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range fields.([]interface{}) {
		if data[k] != nil {

			name := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Name").String()

			result[name] = data[k]
		}
	}

	return result
}

// 获取提交表单的数据
func (p *ResourceImport) getSubmitData(fields interface{}, submitData interface{}) interface{} {

	result := make(map[string]interface{})

	for _, v := range fields.([]interface{}) {

		component := reflect.
			ValueOf(v).
			Elem().
			FieldByName("Component").String()

		name := reflect.
			ValueOf(v).
			Elem().
			FieldByName("Name").String()

		switch component {
		case "inputNumberField":
			result[name] = submitData.(map[string]interface{})[name]
		case "textField":
			result[name] = strings.Trim(submitData.(map[string]interface{})[name].(string), "\n")
		case "selectField":
			options := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Options").Interface()

			result[name] = p.getOptionValue(options, submitData.(map[string]interface{})[name].(string))
		case "cascaderField":
			options := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Options").Interface()

			result[name] = p.getOptionValue(options, submitData.(map[string]interface{})[name].(string))
		case "checkboxField":
			options := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Options").Interface()

			result[name] = p.getOptionValue(options, submitData.(map[string]interface{})[name].(string))
		case "radioField":
			options := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Options").Interface()

			result[name] = p.getOptionValue(options, submitData.(map[string]interface{})[name].(string))
		case "switchField":
			options := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Options").Interface()

			result[name] = p.getSwitchValue(options, submitData.(map[string]interface{})[name].(string))
		default:
			result[name] = submitData.(map[string]interface{})[name]
		}

		if getValue, ok := result[name].([]interface{}); ok {
			result[name], _ = json.Marshal(getValue)
		}

		if getValue, ok := result[name].([]map[string]interface{}); ok {
			result[name], _ = json.Marshal(getValue)
		}

		if getValue, ok := result[name].(map[string]interface{}); ok {
			result[name], _ = json.Marshal(getValue)
		}
	}

	return result
}

// 获取属性值
func (p *ResourceImport) getOptionValue(options interface{}, label string) interface{} {
	var result1 []interface{}
	var result2 interface{}

	labels1 := strings.Split(label, ",")
	labels2 := strings.Split(label, "，")

	if len(labels1) > 1 || len(labels2) > 1 {
		labels := []string{}
		if len(labels1) > 1 {
			labels = labels1
		}

		if len(labels2) > 1 {
			labels = labels2
		}

		for _, v := range options.([]map[string]interface{}) {
			for _, label := range labels {
				if v["label"] == label {
					result1 = append(result1, v["value"])
				}
			}
		}
	} else {
		for _, v := range options.([]map[string]interface{}) {
			if v["label"] == label {
				result2 = v["value"]
			}
		}
	}

	if len(result1) > 0 {
		return result1
	}

	return result2
}

// 获取开关组件值
func (p *ResourceImport) getSwitchValue(options interface{}, label string) interface{} {
	return (options.(map[string]interface{})["on"] == label)
}
