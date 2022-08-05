package requests

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/gobeam/stringy"
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"gorm.io/gorm"
)

type ResourceUpdate struct {
	Quark
}

// 执行行为
func (p *ResourceUpdate) HandleUpdate(c *fiber.Ctx) interface{} {
	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)

	// 获取字段
	fields := resourceInstance.(interface {
		UpdateFields(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).UpdateFields(c, resourceInstance)

	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)

	if data["id"] == "" {
		return msg.Error("参数错误！", "")
	}

	getData := resourceInstance.(interface {
		BeforeSaving(c *fiber.Ctx, data map[string]interface{}) interface{}
	}).BeforeSaving(c, data)

	if value, ok := getData.(error); ok {
		return value
	}

	if value, ok := getData.(map[string]interface{}); ok {
		data = value
	}

	validator := resourceInstance.(interface {
		ValidatorForUpdate(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) error
	}).ValidatorForUpdate(c, resourceInstance, data)

	if validator != nil {
		return validator
	}

	modelInstance := reflect.
		ValueOf(resourceInstance).
		Elem().
		FieldByName("Model").Interface()

	zeroValues := map[string]interface{}{}

	for _, v := range fields.([]interface{}) {

		name := reflect.
			ValueOf(v).
			Elem().
			FieldByName("Name").String()

		formValue := data[name]

		if getValue, ok := formValue.([]interface{}); ok {
			formValue, _ = json.Marshal(getValue)
		}

		if getValue, ok := formValue.([]map[string]interface{}); ok {
			formValue, _ = json.Marshal(getValue)
		}

		if getValue, ok := formValue.(map[string]interface{}); ok {
			formValue, _ = json.Marshal(getValue)
		}

		if name != "" && formValue != nil {
			fieldName := stringy.New(name).CamelCase("?", "")

			reflectFieldName := reflect.
				ValueOf(modelInstance).
				Elem().
				FieldByName(fieldName)

			var reflectValue reflect.Value

			switch reflectFieldName.Type().String() {
			case "int":
				if value, ok := formValue.(bool); ok {
					if value == true {
						reflectValue = reflect.ValueOf(1)
					} else {
						reflectValue = reflect.ValueOf(0)
						zeroValues[fieldName] = 0
					}
				}

				if value, ok := formValue.(float64); ok {
					reflectValue = reflect.ValueOf(int(value))

					// 因为gorm结构体不更新零值，使用map记录零值，单独更新
					if int(value) == 0 {
						zeroValues[fieldName] = 0
					}
				}
			case "time.Time":
				getTime, _ := time.ParseInLocation("2006-01-02 15:04:05", formValue.(string), time.Local)
				reflectValue = reflect.ValueOf(getTime)
			default:
				reflectValue = reflect.ValueOf(formValue)

				if reflect.ValueOf(formValue).Type().String() == "[]uint8" {
					reflectValue = reflect.ValueOf(string(formValue.([]uint8)))
				}
			}

			reflectFieldName.Set(reflectValue)
		}
	}

	// 获取对象
	getModel := model.Where("id = ?", data["id"]).Updates(modelInstance)

	// 因为gorm使用结构体，不更新零值，需要使用map更新零值
	if len(zeroValues) > 0 {
		getModel = model.Where("id = ?", data["id"]).Updates(zeroValues)
	}

	return resourceInstance.(interface {
		AfterSaved(c *fiber.Ctx, model *gorm.DB) interface{}
	}).AfterSaved(c, getModel)
}
