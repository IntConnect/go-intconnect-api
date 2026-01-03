package system_setting

import (
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/storage"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	systemSettingRepository Repository
	localStorageService     *storage.Manager
	auditLogService         auditLog.Service
	validatorService        validator.Service
	dbConnection            *gorm.DB
	viperConfig             *viper.Viper
}

func NewService(systemSettingRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, auditLogService auditLog.Service,
	localStorageService *storage.Manager,

) *ServiceImpl {
	return &ServiceImpl{
		systemSettingRepository: systemSettingRepository,
		auditLogService:         auditLogService,
		validatorService:        validatorService,
		dbConnection:            dbConnection,
		viperConfig:             viperConfig,
		localStorageService:     localStorageService,
	}
}

func (systemSettingService *ServiceImpl) FindAll() []*model.SystemSettingResponse {
	var systemSettingResponsesRequest []*model.SystemSettingResponse
	err := systemSettingService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		systemSettingEntities, err := systemSettingService.systemSettingRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		systemSettingResponsesRequest = helper.MapEntitiesIntoResponses[*entity.SystemSetting, *model.SystemSettingResponse](systemSettingEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return systemSettingResponsesRequest
}

func (systemSettingService *ServiceImpl) FindByKey(systemSettingKey string, isMinimal bool) *model.SystemSettingResponse {
	var systemSettingResponsesRequest *model.SystemSettingResponse

	err := systemSettingService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		systemSettingEntities, err := systemSettingService.systemSettingRepository.FindByKey(gormTransaction, systemSettingKey)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		if isMinimal {
			systemSettingResponsesRequest = helper.MapEntityIntoResponse[*entity.SystemSetting, *model.SystemSettingResponse](systemSettingEntities)
		} else {
			systemSettingResponsesRequest = helper.MapEntityIntoResponseWithIgnoredFields[*entity.SystemSetting, *model.SystemSettingResponse](systemSettingEntities, []string{"Id", "Description", "Key"})
		}

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return systemSettingResponsesRequest
}

// Create - Membuat systemSetting baru
func (systemSettingService *ServiceImpl) Manage(ginContext *gin.Context, createSystemSettingRequest *model.ManageSystemSettingRequest) []*model.SystemSettingResponse {
	valErr := systemSettingService.validatorService.ValidateStruct(createSystemSettingRequest)
	systemSettingService.validatorService.ParseValidationError(valErr, *createSystemSettingRequest)
	err := systemSettingService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		resolvedSchema, isExists := model.SystemSettingSchemas[createSystemSettingRequest.Key]
		if !isExists {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, exception.ErrSystemSettingKeyNotMatch))
		}
		loadedStruct := resolvedSchema.NewPayload
		systemSettingEntity, err := systemSettingService.systemSettingRepository.FindByKey(gormTransaction, createSystemSettingRequest.Key)
		parsedPayload := helper.ParsingHashMapIntoStruct[*model.DashboardSettingPayload](createSystemSettingRequest.Value, loadedStruct().(*model.DashboardSettingPayload))
		modelFile, _ := ginContext.FormFile("value[model]")
		if modelFile != nil {
			(*parsedPayload).ModelFile = modelFile
			newPath, err := systemSettingService.localStorageService.Disk().Put(modelFile, fmt.Sprintf("system-settings/models/%d-%s", time.Now().UnixNano(), modelFile.Filename))
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrSavingResources))
			createSystemSettingRequest.Value["model_path"] = newPath
		} else {
			var modelFileErr error
			if err == nil {
				modelFileErr = systemSettingService.validatorService.ValidateVar((*parsedPayload).ModelFile, "omitempty")
				createSystemSettingRequest.Value["model_path"] = systemSettingEntity.Value["model_path"]
			} else {
				modelFileErr = systemSettingService.validatorService.ValidateVar((*parsedPayload).ModelFile, "required")
			}
			systemSettingService.validatorService.ParseValidationError(valErr, modelFileErr)
		}

		fmt.Println(*parsedPayload)
		valErr = systemSettingService.validatorService.ValidateStruct(*(parsedPayload))
		systemSettingService.validatorService.ParseValidationError(valErr, *parsedPayload)
		systemSettingEntity = helper.MapCreateRequestIntoEntity[model.ManageSystemSettingRequest, entity.SystemSetting](createSystemSettingRequest)
		systemSettingEntity.Value = createSystemSettingRequest.Value
		err = systemSettingService.systemSettingRepository.Manage(gormTransaction, systemSettingEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return systemSettingService.FindAll()
}
