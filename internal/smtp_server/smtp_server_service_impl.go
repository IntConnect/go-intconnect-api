package smtp_server

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
	smtpServerRepository Repository
	auditLogService      auditLog.Service
	validatorService     validator.Service
	dbConnection         *gorm.DB
	viperConfig          *viper.Viper
}

func NewService(smtpServerRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		smtpServerRepository: smtpServerRepository,
		auditLogService:      auditLogService,
		validatorService:     validatorService,
		dbConnection:         dbConnection,
		viperConfig:          viperConfig,
	}
}

func (smtpServerService *ServiceImpl) FindAll() []*model.SmtpServerResponse {
	var smtpServerResponsesRequest []*model.SmtpServerResponse
	err := smtpServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		smtpServerEntities, err := smtpServerService.smtpServerRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		smtpServerResponsesRequest = helper.MapEntitiesIntoResponses[*entity.SmtpServer, *model.SmtpServerResponse](smtpServerEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return smtpServerResponsesRequest
}

func (smtpServerService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.SmtpServerResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var smtpServerResponses []*model.SmtpServerResponse
	var totalItems int64
	err := smtpServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		smtpServerEntities, total, err := smtpServerService.smtpServerRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		smtpServerResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.SmtpServer, *model.SmtpServerResponse](
			smtpServerEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"MQTT Broker fetched successfully",
		smtpServerResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat smtpServer baru
func (smtpServerService *ServiceImpl) Create(ginContext *gin.Context, createSmtpServerRequest *model.CreateSmtpServerRequest) *model.PaginatedResponse[*model.SmtpServerResponse] {
	jwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	ipAddress, userAgent := helper.ExtractRequestMeta(ginContext)
	var paginationResp *model.PaginatedResponse[*model.SmtpServerResponse]
	valErr := smtpServerService.validatorService.ValidateStruct(createSmtpServerRequest)
	smtpServerService.validatorService.ParseValidationError(valErr, *createSmtpServerRequest)
	err := smtpServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		smtpServerEntity := helper.MapCreateRequestIntoEntity[model.CreateSmtpServerRequest, entity.SmtpServer](createSmtpServerRequest)
		err := smtpServerService.smtpServerRepository.Create(gormTransaction, smtpServerEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		smtpServerService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      jwtClaims.Id,
			Action:      model.AUDIT_LOG_CREATE,
			Feature:     model.AUDIT_LOG_FEATURE_SMTP_SERVER,
			Description: "",
			Before:      nil,
			After:       smtpServerEntity,
			IpAddress:   ipAddress,
			UserAgent:   userAgent,
		})
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = smtpServerService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (smtpServerService *ServiceImpl) Update(ginContext *gin.Context, updateSmtpServerRequest *model.UpdateSmtpServerRequest) *model.PaginatedResponse[*model.SmtpServerResponse] {
	jwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	ipAddress, userAgent := helper.ExtractRequestMeta(ginContext)
	var paginationResp *model.PaginatedResponse[*model.SmtpServerResponse]
	valErr := smtpServerService.validatorService.ValidateStruct(updateSmtpServerRequest)
	smtpServerService.validatorService.ParseValidationError(valErr, *updateSmtpServerRequest)
	err := smtpServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		smtpServer, err := smtpServerService.smtpServerRepository.FindById(gormTransaction, updateSmtpServerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateSmtpServerRequest, smtpServer)
		err = smtpServerService.smtpServerRepository.Update(gormTransaction, smtpServer)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		smtpServerService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      jwtClaims.Id,
			Action:      model.AUDIT_LOG_CREATE,
			Feature:     model.AUDIT_LOG_FEATURE_SMTP_SERVER,
			Description: "",
			Before:      nil,
			After:       smtpServer,
			IpAddress:   ipAddress,
			UserAgent:   userAgent,
		})
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = smtpServerService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (smtpServerService *ServiceImpl) Delete(ginContext *gin.Context, deleteSmtpServerRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.SmtpServerResponse] {
	jwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	ipAddress, userAgent := helper.ExtractRequestMeta(ginContext)
	var paginationResp *model.PaginatedResponse[*model.SmtpServerResponse]
	valErr := smtpServerService.validatorService.ValidateStruct(deleteSmtpServerRequest)
	smtpServerService.validatorService.ParseValidationError(valErr, *deleteSmtpServerRequest)
	err := smtpServerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		smtpServerEntity, err := smtpServerService.smtpServerRepository.FindById(gormTransaction, deleteSmtpServerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = smtpServerService.smtpServerRepository.Delete(gormTransaction, deleteSmtpServerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		paginationResp = smtpServerService.FindAllPagination(&paginationRequest)
		smtpServerService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      jwtClaims.Id,
			Action:      model.AUDIT_LOG_DELETE,
			Feature:     model.AUDIT_LOG_FEATURE_SMTP_SERVER,
			Description: deleteSmtpServerRequest.Reason,
			Before:      smtpServerEntity,
			After:       "",
			IpAddress:   ipAddress,
			UserAgent:   userAgent,
		})
		return nil
	})
	paginationRequest := model.NewPaginationRequest()
	paginationResp = smtpServerService.FindAllPagination(&paginationRequest)

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}
