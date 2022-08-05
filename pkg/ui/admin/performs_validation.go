package admin

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/derekstavis/go-qs"
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/db"
)

// 创建请求的验证器
func (p *Resource) ValidatorForCreation(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) error {
	rules, messages := p.RulesForCreation(c, resourceInstance)

	validator := p.Validator(rules, messages, data)

	p.afterValidation(c, validator)
	p.afterCreationValidation(c, validator)

	return validator
}

// 验证规则
func (p *Resource) Validator(rules []interface{}, messages []interface{}, data map[string]interface{}) error {
	var result error

	for _, rule := range rules {
		for k, v := range rule.(map[string]interface{}) {
			fieldValue := data[k]
			for _, item := range v.([]interface{}) {
				getItem, ok := item.(string)
				if ok {
					getItems := strings.Split(getItem, ":")
					getOption := ""
					if len(getItems) == 2 {
						getItem = getItems[0]
						getOption = getItems[1]
					}

					switch getItem {
					case "required":
						if fieldValue == nil {
							errMsg := p.getRuleMessage(messages, k+"."+getItem)

							if errMsg != "" {
								result = errors.New(errMsg)
							}
						}
					case "min":
						strNum := utf8.RuneCountInString(fieldValue.(string))
						minOption, _ := strconv.Atoi(getOption)

						if strNum < minOption {
							errMsg := p.getRuleMessage(messages, k+"."+getItem)
							if errMsg != "" {
								result = errors.New(errMsg)
							}
						}
					case "max":
						strNum := utf8.RuneCountInString(fieldValue.(string))
						maxOption, _ := strconv.Atoi(getOption)

						if strNum > maxOption {
							errMsg := p.getRuleMessage(messages, k+"."+getItem)
							if errMsg != "" {
								result = errors.New(errMsg)
							}
						}
					case "unique":
						var (
							table  string
							field  string
							except string
							count  int
						)

						uniqueOptions := strings.Split(getOption, ",")

						if len(uniqueOptions) == 2 {
							table = uniqueOptions[0]
							field = uniqueOptions[1]
						}

						if len(uniqueOptions) == 3 {
							table = uniqueOptions[0]
							field = uniqueOptions[1]
							except = uniqueOptions[2]
						}

						if except != "" {
							(&db.Model{}).DB().Raw("SELECT Count("+field+") FROM "+table+" WHERE id <> "+except+" AND "+field+" = ?", fieldValue).Scan(&count)
						} else {
							(&db.Model{}).DB().Raw("SELECT Count("+field+") FROM "+table+" WHERE "+field+" = ?", fieldValue).Scan(&count)
						}

						if count > 0 {
							errMsg := p.getRuleMessage(messages, k+"."+getItem)
							if errMsg != "" {
								result = errors.New(errMsg)
							}
						}
					}
				}
			}
		}
	}

	return result
}

// 获取规则错误信息
func (p *Resource) getRuleMessage(messages []interface{}, key string) string {
	msg := ""

	for _, v := range messages {
		getMsg := v.(map[string]interface{})[key]
		if getMsg != nil {
			msg = getMsg.(string)
		}
	}

	return msg
}

// 创建请求的验证规则
func (p *Resource) RulesForCreation(c *fiber.Ctx, resourceInstance interface{}) ([]interface{}, []interface{}) {

	fields := resourceInstance.(interface {
		CreationFieldsWithoutWhen(*fiber.Ctx, interface{}) interface{}
	}).CreationFieldsWithoutWhen(c, resourceInstance)

	rules := []interface{}{}
	ruleMessages := []interface{}{}

	for _, v := range fields.([]interface{}) {
		getResult := p.getRulesForCreation(c, v)

		if len(getResult["rules"].(map[string]interface{})) > 0 {
			rules = append(rules, getResult["rules"])
		}

		if len(getResult["messages"].(map[string]interface{})) > 0 {
			ruleMessages = append(ruleMessages, getResult["messages"])
		}

		when := reflect.
			ValueOf(v).
			Elem().
			FieldByName("When").Interface()

		if when != nil {
			whenItems := reflect.
				ValueOf(when).
				Elem().
				FieldByName("Items").Interface()

			if whenItems != nil {
				for _, vi := range whenItems.([]map[string]interface{}) {
					if p.needValidateWhenRules(c, vi) {
						body := vi["body"]
						if body != nil {
							// 如果为数组
							getBody, ok := body.([]interface{})
							if ok {
								for _, bv := range getBody {
									whenItemResult := p.getRulesForCreation(c, bv)

									if len(whenItemResult["rules"].(map[string]interface{})) > 0 {
										rules = append(rules, whenItemResult["rules"])
									}

									if len(whenItemResult["messages"].(map[string]interface{})) > 0 {
										ruleMessages = append(ruleMessages, whenItemResult["messages"])
									}
								}
							} else {
								whenItemResult := p.getRulesForCreation(c, getBody)

								if len(whenItemResult["rules"].(map[string]interface{})) > 0 {
									rules = append(rules, whenItemResult["rules"])
								}

								if len(whenItemResult["messages"].(map[string]interface{})) > 0 {
									ruleMessages = append(ruleMessages, whenItemResult["messages"])
								}
							}
						}
					}
				}
			}
		}

	}

	return rules, ruleMessages
}

// 判断是否需要验证When组件里的规则
func (p *Resource) needValidateWhenRules(c *fiber.Ctx, when map[string]interface{}) bool {
	conditionName := when["condition_name"].(string)
	conditionOption := when["condition_option"]
	conditionOperator := when["condition_operator"].(string)
	result := false

	data, error := qs.Unmarshal(c.OriginalURL())
	if error != nil {
		return false
	}

	value, ok := data[conditionName]
	if !ok {
		return false
	}

	valueString, isString := value.(string)
	if isString {
		if valueString == "" {
			return false
		}
	}

	switch conditionOperator {
	case "=":
		result = (value.(string) == conditionOption.(string))
	case ">":
		result = (value.(string) > conditionOption.(string))
	case "<":
		result = (value.(string) < conditionOption.(string))
	case "<=":
		result = (value.(string) <= conditionOption.(string))
	case ">=":
		result = (value.(string) >= conditionOption.(string))
	case "has":
		_, isArray := value.([]string)
		if isArray {
			getJson, err := json.Marshal(value)
			if err != nil {
				result = strings.Contains(string(getJson), conditionOption.(string))
			}
		} else {
			result = strings.Contains(value.(string), conditionOption.(string))
		}
	case "in":
		conditionOptionArray, isArray := conditionOption.([]string)
		if isArray {
			for _, v := range conditionOptionArray {
				if v == value.(string) {
					result = true
				}
			}
		}
	default:
		result = (value.(string) == conditionOption)
	}

	return result
}

// 获取创建请求资源规则
func (p *Resource) getRulesForCreation(c *fiber.Ctx, field interface{}) map[string]interface{} {
	getRules := map[string]interface{}{}
	getRuleMessages := map[string]interface{}{}

	name := reflect.
		ValueOf(field).
		Elem().
		FieldByName("Name").String()

	rules := reflect.
		ValueOf(field).
		Elem().
		FieldByName("Rules").Interface()

	ruleMessages := reflect.
		ValueOf(field).
		Elem().
		FieldByName("RuleMessages").Interface()

	creationRules := reflect.
		ValueOf(field).
		Elem().
		FieldByName("CreationRules").Interface()

	creationRuleMessages := reflect.
		ValueOf(field).
		Elem().
		FieldByName("CreationRuleMessages").Interface()

	items := []interface{}{}

	for _, v := range p.formatRules(c, rules.([]string)) {
		items = append(items, v)
	}

	for key, v := range ruleMessages.(map[string]string) {
		getRuleMessages[name+"."+key] = v
	}

	for _, v := range p.formatRules(c, creationRules.([]string)) {
		items = append(items, v)
	}

	for key, v := range creationRuleMessages.(map[string]string) {
		getRuleMessages[name+"."+key] = v
	}

	if len(items) > 0 {
		getRules[name] = items
	}

	return map[string]interface{}{
		"rules":    getRules,
		"messages": getRuleMessages,
	}
}

// 更新请求的验证器
func (p *Resource) ValidatorForUpdate(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) error {
	rules, messages := p.RulesForUpdate(c, resourceInstance)

	validator := p.Validator(rules, messages, data)

	p.afterValidation(c, validator)
	p.afterCreationValidation(c, validator)

	return validator
}

// 更新请求的验证规则
func (p *Resource) RulesForUpdate(c *fiber.Ctx, resourceInstance interface{}) ([]interface{}, []interface{}) {

	fields := resourceInstance.(interface {
		UpdateFieldsWithoutWhen(*fiber.Ctx, interface{}) interface{}
	}).UpdateFieldsWithoutWhen(c, resourceInstance)

	rules := []interface{}{}
	ruleMessages := []interface{}{}

	for _, v := range fields.([]interface{}) {
		getResult := p.getRulesForUpdate(c, v)

		if len(getResult["rules"].(map[string]interface{})) > 0 {
			rules = append(rules, getResult["rules"])
		}

		if len(getResult["messages"].(map[string]interface{})) > 0 {
			ruleMessages = append(ruleMessages, getResult["messages"])
		}

		when := reflect.
			ValueOf(v).
			Elem().
			FieldByName("When").Interface()

		if when != nil {
			whenItems := reflect.
				ValueOf(when).
				Elem().
				FieldByName("Items").Interface()

			if whenItems != nil {
				for _, vi := range whenItems.([]map[string]interface{}) {
					if p.needValidateWhenRules(c, vi) {
						body := vi["body"]

						if body != nil {

							// 如果为数组
							getBody, ok := body.([]interface{})
							if ok {
								for _, bv := range getBody {
									whenItemResult := p.getRulesForUpdate(c, bv)

									if len(whenItemResult["rules"].(map[string]interface{})) > 0 {
										rules = append(rules, whenItemResult["rules"])
									}

									if len(whenItemResult["messages"].(map[string]interface{})) > 0 {
										ruleMessages = append(ruleMessages, whenItemResult["messages"])
									}
								}
							} else {
								whenItemResult := p.getRulesForUpdate(c, getBody)

								if len(whenItemResult["rules"].(map[string]interface{})) > 0 {
									rules = append(rules, whenItemResult["rules"])
								}

								if len(whenItemResult["messages"].(map[string]interface{})) > 0 {
									ruleMessages = append(ruleMessages, whenItemResult["messages"])
								}
							}
						}
					}
				}
			}
		}

	}

	return rules, ruleMessages
}

// 获取更新请求资源规则
func (p *Resource) getRulesForUpdate(c *fiber.Ctx, field interface{}) map[string]interface{} {

	getRules := map[string]interface{}{}
	getRuleMessages := map[string]interface{}{}

	name := reflect.
		ValueOf(field).
		Elem().
		FieldByName("Name").String()

	rules := reflect.
		ValueOf(field).
		Elem().
		FieldByName("Rules").Interface()

	ruleMessages := reflect.
		ValueOf(field).
		Elem().
		FieldByName("RuleMessages").Interface()

	updateRules := reflect.
		ValueOf(field).
		Elem().
		FieldByName("UpdateRules").Interface()

	updateRuleMessages := reflect.
		ValueOf(field).
		Elem().
		FieldByName("UpdateRuleMessages").Interface()

	items := []interface{}{}

	for _, v := range p.formatRules(c, rules.([]string)) {
		items = append(items, v)
	}

	for key, v := range ruleMessages.(map[string]string) {
		getRuleMessages[name+"."+key] = v
	}

	for _, v := range p.formatRules(c, updateRules.([]string)) {
		items = append(items, v)
	}

	for key, v := range updateRuleMessages.(map[string]string) {
		getRuleMessages[name+"."+key] = v
	}

	if len(items) > 0 {
		getRules[name] = items
	}

	return map[string]interface{}{
		"rules":    getRules,
		"messages": getRuleMessages,
	}
}

// 导入请求的验证器
func (p *Resource) ValidatorForImport(c *fiber.Ctx, resourceInstance interface{}, data map[string]interface{}) error {
	rules, messages := p.RulesForImport(c, resourceInstance)

	validator := p.Validator(rules, messages, data)

	p.afterValidation(c, validator)
	p.afterCreationValidation(c, validator)

	return validator
}

// 创建请求的验证规则
func (p *Resource) RulesForImport(c *fiber.Ctx, resourceInstance interface{}) ([]interface{}, []interface{}) {

	fields := resourceInstance.(interface {
		ImportFieldsWithoutWhen(*fiber.Ctx, interface{}) interface{}
	}).ImportFieldsWithoutWhen(c, resourceInstance)

	rules := []interface{}{}
	ruleMessages := []interface{}{}

	for _, v := range fields.([]interface{}) {
		getResult := p.getRulesForCreation(c, v)

		if len(getResult["rules"].(map[string]interface{})) > 0 {
			rules = append(rules, getResult["rules"])
		}

		if len(getResult["messages"].(map[string]interface{})) > 0 {
			ruleMessages = append(ruleMessages, getResult["messages"])
		}

		when := reflect.
			ValueOf(v).
			Elem().
			FieldByName("When").Interface()

		if when != nil {
			whenItems := reflect.
				ValueOf(when).
				Elem().
				FieldByName("Items").Interface()

			if whenItems != nil {
				for _, vi := range whenItems.([]map[string]interface{}) {
					if p.needValidateWhenRules(c, vi) {
						body := vi["body"]

						if body != nil {

							// 如果为数组
							getBody, ok := body.([]interface{})
							if ok {
								for _, bv := range getBody {
									whenItemResult := p.getRulesForCreation(c, bv)

									if len(whenItemResult["rules"].(map[string]interface{})) > 0 {
										rules = append(rules, whenItemResult["rules"])
									}

									if len(whenItemResult["messages"].(map[string]interface{})) > 0 {
										ruleMessages = append(ruleMessages, whenItemResult["messages"])
									}
								}
							} else {
								whenItemResult := p.getRulesForCreation(c, getBody)

								if len(whenItemResult["rules"].(map[string]interface{})) > 0 {
									rules = append(rules, whenItemResult["rules"])
								}

								if len(whenItemResult["messages"].(map[string]interface{})) > 0 {
									ruleMessages = append(ruleMessages, whenItemResult["messages"])
								}
							}
						}
					}
				}
			}
		}

	}

	return rules, ruleMessages
}

// 格式化规则
func (p *Resource) formatRules(c *fiber.Ctx, rules []string) []string {
	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)

	formId := data["id"]
	requestId := c.Query("id")

	if requestId == "" && formId == nil {
		return rules
	}

	if requestId != "" {
		for key, v := range rules {
			rules[key] = strings.Replace(v, "{id}", requestId, -1)
		}
	} else if formId != nil {
		for key, v := range rules {
			requestId = strconv.FormatFloat(formId.(float64), 'E', -1, 64)
			rules[key] = strings.Replace(v, "{id}", requestId, -1)
		}
	}

	return rules
}

// 验证完成后回调
func (p *Resource) afterValidation(c *fiber.Ctx, validator interface{}) {
	//
}

// 创建请求验证完成后回调
func (p *Resource) afterCreationValidation(c *fiber.Ctx, validator interface{}) {
	//
}

// 更新请求验证完成后回调
func (p *Resource) afterUpdateValidation(c *fiber.Ctx, validator interface{}) {
	//
}

// 创建请求验证完成后回调
func (p *Resource) afterImportValidation(c *fiber.Ctx, validator interface{}) {
	//
}
