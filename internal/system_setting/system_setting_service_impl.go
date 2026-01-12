package system_setting

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/storage"
	"go-intconnect-api/internal/system_setting/processor"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"

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
	systemSettingRegistry   *Registry
}

func NewService(systemSettingRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, auditLogService auditLog.Service,
	localStorageService *storage.Manager,
	systemSettingRegistry *Registry,
) *ServiceImpl {
	systemSettingRegistry.Register(
		"DASHBOARD_SETTINGS",
		processor.NewDashboardSettingHandler(localStorageService, validatorService),
	)
	systemSettingRegistry.Register(
		"LISTENER_SETTINGS",
		processor.NewListenerSettingHandler(validatorService),
	)
	return &ServiceImpl{
		systemSettingRepository: systemSettingRepository,
		auditLogService:         auditLogService,
		validatorService:        validatorService,
		dbConnection:            dbConnection,
		viperConfig:             viperConfig,
		localStorageService:     localStorageService,
		systemSettingRegistry:   systemSettingRegistry,
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
		systemSettingEntities, _ := systemSettingService.systemSettingRepository.FindByKey(gormTransaction, systemSettingKey)

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

func (systemSettingService *ServiceImpl) Manage(
	ginContext *gin.Context,
	manageSystemSettingRequest *model.ManageSystemSettingRequest,
) []*model.SystemSettingResponse {

	valErr := systemSettingService.validatorService.ValidateStruct(manageSystemSettingRequest)
	systemSettingService.validatorService.ParseValidationError(valErr, *manageSystemSettingRequest)

	err := systemSettingService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {

		existingSystemSettingEntity, err := systemSettingService.systemSettingRepository.FindByKey(gormTransaction, manageSystemSettingRequest.Key)
		systemSettingProcessor := systemSettingService.systemSettingRegistry.Resolve(manageSystemSettingRequest.Key)
		systemSettingEntity, err := systemSettingProcessor.Handle(ginContext, gormTransaction, existingSystemSettingEntity, manageSystemSettingRequest)
		if err != nil {
			return err
		}

		if err := systemSettingService.systemSettingRepository.Manage(gormTransaction, systemSettingEntity); err != nil {
			return exception.ParseGormError(err)
		}

		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return systemSettingService.FindAll()
}
