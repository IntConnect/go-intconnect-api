package processor

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListenerSettingHandler struct {
	validatorService validator.Service
}

func NewListenerSettingHandler(
	validatorService validator.Service,
) *ListenerSettingHandler {
	return &ListenerSettingHandler{
		validatorService: validatorService,
	}
}

func (listenerSettingHandler *ListenerSettingHandler) Handle(
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
	parsedPayload := helper.ParsingHashMapIntoStruct[*model.ListenerSettingPayload](manageSystemSettingRequest.Value, loadedStruct().(*model.ListenerSettingPayload))

	err := listenerSettingHandler.validatorService.ValidateStruct(*(parsedPayload))
	listenerSettingHandler.validatorService.ParseValidationError(err, *parsedPayload)
	updatedSystemSettingEntity := helper.MapCreateRequestIntoEntity[
		model.ManageSystemSettingRequest,
		entity.SystemSetting,
	](manageSystemSettingRequest)

	systemSettingEntity.Value = updatedSystemSettingEntity.Value
	return updatedSystemSettingEntity, nil
}
