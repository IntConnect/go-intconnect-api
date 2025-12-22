package system_setting

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllSystemSetting(ginContext *gin.Context)
	FindSystemSettingByKey(ginContext *gin.Context)
	ManageSystemSetting(ginContext *gin.Context)
	FindMinimalSystemSettingByKey(ginContext *gin.Context)
}
