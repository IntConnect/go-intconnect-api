package system_setting

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemSettingHandler interface {
	Handle(
		ginContext *gin.Context,
		gormTransaction *gorm.DB,
		systemSettingEntity *entity.SystemSetting,
		manageSystemSettingRequest *model.ManageSystemSettingRequest,
	) (*entity.SystemSetting, error)
}
