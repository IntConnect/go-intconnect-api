package modbus_server

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	modbusServerRepository Repository
	auditLogService        auditLog.Service
	validatorService       validator.Service
	dbConnection           *gorm.DB
	viperConfig            *viper.Viper
}

func NewService(modbusServerRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		modbusServerRepository: modbusServerRepository,
		validatorService:       validatorService,
		dbConnection:           dbConnection,
		viperConfig:            viperConfig,
		auditLogService:        auditLogService,
	}
}

// Create - Membuat modbusServer baru
func (modbusServerService *ServiceImpl) FindAll() []*model.ModbusServerResponse {
	var allModbusServer []*model.ModbusServerResponse
	err := modbusServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modbusServerResponse, err := modbusServerService.modbusServerRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allModbusServer = helper.MapEntitiesIntoResponses[entity.ModbusServer, *model.ModbusServerResponse](modbusServerResponse)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allModbusServer
}

// Create - Membuat modbusServer baru
func (modbusServerService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.ModbusServerResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var modbusServerResponses []*model.ModbusServerResponse
	var totalItems int64

	err := modbusServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modbusServerEntities, total, err := modbusServerService.modbusServerRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		modbusServerResponses = helper.MapEntitiesIntoResponsesWithFunc[entity.ModbusServer, *model.ModbusServerResponse](
			modbusServerEntities,
			mapper.FuncMapAuditable,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Mqtt topic fetched successfully",
		modbusServerResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat modbusServer baru
func (modbusServerService *ServiceImpl) Create(ginContext *gin.Context, createModbusServerRequest *model.CreateModbusServerRequest) *model.PaginatedResponse[*model.ModbusServerResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginatedResp *model.PaginatedResponse[*model.ModbusServerResponse]
	valErr := modbusServerService.validatorService.ValidateStruct(createModbusServerRequest)
	modbusServerService.validatorService.ParseValidationError(valErr, *createModbusServerRequest)
	err := modbusServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modbusServerEntity := helper.MapCreateRequestIntoEntity[model.CreateModbusServerRequest, entity.ModbusServer](createModbusServerRequest)
		modbusServerEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		err := modbusServerService.modbusServerRepository.Create(gormTransaction, modbusServerEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		_ = modbusServerService.auditLogService.Record(
			ginContext,
			model.AUDIT_LOG_UPDATE,
			model.AUDIT_LOG_FEATURE_MODBUS_SERVER,
			nil,
			modbusServerEntity,
			"",
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = modbusServerService.FindAllPagination(&paginationRequest)
	return paginatedResp
}

func (modbusServerService *ServiceImpl) Update(ginContext *gin.Context, updateModbusServerRequest *model.UpdateModbusServerRequest) *model.PaginatedResponse[*model.ModbusServerResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginatedResp *model.PaginatedResponse[*model.ModbusServerResponse]
	valErr := modbusServerService.validatorService.ValidateStruct(updateModbusServerRequest)
	modbusServerService.validatorService.ParseValidationError(valErr, *updateModbusServerRequest)
	err := modbusServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modbusServer, err := modbusServerService.modbusServerRepository.FindById(gormTransaction, updateModbusServerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		rawModbusServer := modbusServer
		helper.MapUpdateRequestIntoEntity(updateModbusServerRequest, modbusServer)
		modbusServer.Auditable = entity.UpdateAuditable(userJwtClaims.Username)
		err = modbusServerService.modbusServerRepository.Update(gormTransaction, modbusServer)
		_ = modbusServerService.auditLogService.Record(
			ginContext,
			model.AUDIT_LOG_UPDATE,
			model.AUDIT_LOG_FEATURE_MODBUS_SERVER,
			rawModbusServer,
			modbusServer,
			"",
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = modbusServerService.FindAllPagination(&paginationRequest)
	return paginatedResp
}

func (modbusServerService *ServiceImpl) Delete(ginContext *gin.Context, deleteModbusServerRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.ModbusServerResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginatedResp *model.PaginatedResponse[*model.ModbusServerResponse]
	valErr := modbusServerService.validatorService.ValidateStruct(deleteModbusServerRequest)
	modbusServerService.validatorService.ParseValidationError(valErr, *deleteModbusServerRequest)
	err := modbusServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modbusServer, err := modbusServerService.modbusServerRepository.FindById(gormTransaction, deleteModbusServerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		modbusServer.Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		err = modbusServerService.modbusServerRepository.Delete(gormTransaction, deleteModbusServerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		_ = modbusServerService.auditLogService.Record(
			ginContext,
			model.AUDIT_LOG_CREATE,
			model.AUDIT_LOG_FEATURE_MODBUS_SERVER,
			nil,
			modbusServer,
			deleteModbusServerRequest.Reason,
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = modbusServerService.FindAllPagination(&paginationRequest)
	return paginatedResp
}
