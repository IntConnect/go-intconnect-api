package processor

import (
	"encoding/json"
	"fmt"
	"go-intconnect-api/configs"
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
	redisInstance    *configs.RedisInstance
}

func NewListenerSettingHandler(
	validatorService validator.Service,
	redisInstance *configs.RedisInstance,
) *ListenerSettingHandler {
	return &ListenerSettingHandler{
		validatorService: validatorService,
		redisInstance:    redisInstance,
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
	marshalledListenerSettingEntity, err := json.Marshal(updatedSystemSettingEntity.Value)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError))
	fmt.Println(marshalledListenerSettingEntity, updatedSystemSettingEntity.Value)
	redisKey := fmt.Sprintf("dashboard_settings:listener_setting")
	err = listenerSettingHandler.redisInstance.RedisClient.Set(
		ginContext.Request.Context(),
		redisKey,
		marshalledListenerSettingEntity,
		0,
	).Err()
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError))

	return updatedSystemSettingEntity, nil
}
