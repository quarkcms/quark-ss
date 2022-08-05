package requests

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/gobeam/stringy"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ResourceStore struct {
	Quark
}

// 执行行为
func (p *ResourceStore) HandleStore(c *fiber.Ctx) interface{} {
	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)

	// 获取字段
	fields := resourceInstance.(interface {
		CreationFields(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).CreationFields(c, resourceInstance)

	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)

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
		ValidatorForCreation(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) error
	}).ValidatorForCreation(c, resourceInstance, data)

	if validator != nil {
		return validator
	}

	modelInstance := reflect.
		ValueOf(resourceInstance).
		Elem().
		FieldByName("Model").Interface()

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
					}
				}

				if value, ok := formValue.(float64); ok {
					reflectValue = reflect.ValueOf(int(value))
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
	getModel := model.Create(modelInstance)

	return resourceInstance.(interface {
		AfterSaved(c *fiber.Ctx, model *gorm.DB) interface{}
	}).AfterSaved(c, getModel)
}
