package system_setting

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.SystemSettingResponse
	FindByKey(systemSettingKey string, isMinimal bool) *model.SystemSettingResponse
	Manage(ginContext *gin.Context, createSystemSettingRequest *model.ManageSystemSettingRequest) []*model.SystemSettingResponse
}
