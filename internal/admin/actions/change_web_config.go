package actions

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"gorm.io/gorm"
)

type ChangeWebConfig struct {
	actions.Action
}

// 执行行为句柄
func (p *ChangeWebConfig) Handle(c *fiber.Ctx, model *gorm.DB) error {

	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)
	result := true

	for k, v := range data {
		config := map[string]interface{}{}
		(&db.Model{}).Model(&models.Config{}).Where("name =?", k).First(&config)

		if getValue, ok := v.([]interface{}); ok {
			v, _ = json.Marshal(getValue)
		}

		if getValue, ok := v.([]map[string]interface{}); ok {
			v, _ = json.Marshal(getValue)
		}

		if getValue, ok := v.(map[string]interface{}); ok {
			v, _ = json.Marshal(getValue)
		}

		updateResult := (&db.Model{}).Model(&models.Config{}).Where("name", k).Update("value", v)

		if updateResult.Error != nil {
			result = false
		}
	}

	if result {

		// 刷新网站配置
		utils.RefreshWebConfig()

		// 返回成功
		return msg.Success("操作成功！", "", "")
	} else {
		return msg.Error("操作失败，请重试！", "")
	}
}
