package system_setting

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/spf13/viper"
)

type Handler struct {
	systemSettingService Service
	formDecoder          *form.Decoder
	viperConfig          *viper.Viper
}

func NewHandler(systemSettingService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		systemSettingService: systemSettingService,
		viperConfig:          viperConfig,
		formDecoder:          form.NewDecoder(),
	}
}

func (systemSettingHandler *Handler) FindAllSystemSetting(ginContext *gin.Context) {
	systemSettingResponses := systemSettingHandler.systemSettingService.FindAll()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("System Setting has been fetched", systemSettingResponses))
}

func (systemSettingHandler *Handler) FindSystemSettingByKey(ginContext *gin.Context) {
	systemSettingKey := ginContext.Param("key")
	systemSettingResponse := systemSettingHandler.systemSettingService.FindByKey(systemSettingKey, false)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("System Setting has been fetched", systemSettingResponse))
}

func (systemSettingHandler *Handler) FindMinimalSystemSettingByKey(ginContext *gin.Context) {
	systemSettingKey := ginContext.Param("key")
	systemSettingResponse := systemSettingHandler.systemSettingService.FindByKey(systemSettingKey, true)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("System Setting has been fetched", systemSettingResponse))
}

func (systemSettingHandler *Handler) ManageSystemSetting(ginContext *gin.Context) {
	var createSystemSettingModel model.ManageSystemSettingRequest
	err := ginContext.Request.ParseMultipartForm(500 << 20) // 32MB maxMemory
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	err = systemSettingHandler.formDecoder.Decode(&createSystemSettingModel, ginContext.Request.PostForm)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := systemSettingHandler.systemSettingService.Manage(ginContext, &createSystemSettingModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
