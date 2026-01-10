package processor

import (
	"fmt"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/storage"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardSettingHandler struct {
	localStorage     *storage.Manager
	validatorService validator.Service
}

func NewDashboardSettingHandler(
	localStorage *storage.Manager,
	validatorService validator.Service,
) *DashboardSettingHandler {
	return &DashboardSettingHandler{
		localStorage:     localStorage,
		validatorService: validatorService,
	}
}

func (dashboardSettingHandler *DashboardSettingHandler) Handle(
	ginContext *gin.Context,
	gormTransaction *gorm.DB,
	systemSettingEntity *entity.SystemSetting,
	manageSystemSettingRequest *model.ManageSystemSettingRequest,
) (*entity.SystemSetting, error) {
	resolvedSchema, isExists := model.SystemSettingSchemas[manageSystemSettingRequest.Key]
	if !isExists {
		exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, exception.ErrSystemSettingKeyNotMatch))
	}
	loadedStruct := resolvedSchema.NewPayload
	parsedPayload := helper.ParsingHashMapIntoStruct[*model.DashboardSettingPayload](manageSystemSettingRequest.Value, loadedStruct().(*model.DashboardSettingPayload))

	modelFile, _ := ginContext.FormFile("value[model]")
	if modelFile != nil {
		(*parsedPayload).ModelFile = modelFile
		path, err := dashboardSettingHandler.localStorage.Disk().Put(
			modelFile,
			fmt.Sprintf("system-settings/models/%d-%s", time.Now().UnixNano(), modelFile.Filename),
		)
		if err != nil {
			return nil, exception.NewApplicationError(400, exception.ErrSavingResources)
		}
		manageSystemSettingRequest.Value["model_path"] = path
	} else if systemSettingEntity != nil {
		manageSystemSettingRequest.Value["model_path"] = systemSettingEntity.Value["model_path"]
	}

	updatedSystemSettingEntity := helper.MapCreateRequestIntoEntity[
		model.ManageSystemSettingRequest,
		entity.SystemSetting,
	](manageSystemSettingRequest)

	updatedSystemSettingEntity.Value = updatedSystemSettingEntity.Value
	return updatedSystemSettingEntity, nil
}
