package breakdown

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
	breakdownRepository Repository
	auditLogService     auditLog.Service
	validatorService    validator.Service
	dbConnection        *gorm.DB
	viperConfig         *viper.Viper
}

func NewService(breakdownRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		breakdownRepository: breakdownRepository,
		validatorService:    validatorService,
		dbConnection:        dbConnection,
		viperConfig:         viperConfig,
		auditLogService:     auditLogService,
	}
}

// Create - Membuat breakdown baru
func (breakdownService *ServiceImpl) FindAll() []*model.BreakdownResponse {
	var allBreakdown []*model.BreakdownResponse
	err := breakdownService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		breakdownResponse, err := breakdownService.breakdownRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allBreakdown = helper.MapEntitiesIntoResponses[entity.Breakdown, *model.BreakdownResponse](breakdownResponse)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allBreakdown
}

// Create - Membuat breakdown baru
func (breakdownService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.BreakdownResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var breakdownResponses []*model.BreakdownResponse
	var totalItems int64

	err := breakdownService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		breakdownEntities, total, err := breakdownService.breakdownRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		breakdownResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.Breakdown, *model.BreakdownResponse](
			breakdownEntities,
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
		breakdownResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat breakdown baru
func (breakdownService *ServiceImpl) Create(ginContext *gin.Context, createBreakdownRequest *model.CreateBreakdownRequest) *model.PaginatedResponse[*model.BreakdownResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginatedResp *model.PaginatedResponse[*model.BreakdownResponse]
	valErr := breakdownService.validatorService.ValidateStruct(createBreakdownRequest)
	breakdownService.validatorService.ParseValidationError(valErr, *createBreakdownRequest)
	err := breakdownService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		breakdownEntity := helper.MapCreateRequestIntoEntity[model.CreateBreakdownRequest, entity.Breakdown](createBreakdownRequest)
		breakdownEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		err := breakdownService.breakdownRepository.Create(gormTransaction, breakdownEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := breakdownService.auditLogService.Build(
			nil,             // before entity
			breakdownEntity, // after entity
			nil,
			"",
		)

		err = breakdownService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_BREAKDOWN,
				auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = breakdownService.FindAllPagination(&paginationRequest)
	return paginatedResp
}

func (breakdownService *ServiceImpl) Update(ginContext *gin.Context, updateBreakdownRequest *model.UpdateBreakdownRequest) *model.PaginatedResponse[*model.BreakdownResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginatedResp *model.PaginatedResponse[*model.BreakdownResponse]
	valErr := breakdownService.validatorService.ValidateStruct(updateBreakdownRequest)
	breakdownService.validatorService.ParseValidationError(valErr, *updateBreakdownRequest)
	err := breakdownService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		breakdownEntity, err := breakdownService.breakdownRepository.FindById(gormTransaction, updateBreakdownRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		beforeBreakdown := *breakdownEntity
		helper.MapUpdateRequestIntoEntity(updateBreakdownRequest, breakdownEntity)
		breakdownEntity.Auditable = entity.UpdateAuditable(userJwtClaims.Username)
		err = breakdownService.breakdownRepository.Update(gormTransaction, breakdownEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := breakdownService.auditLogService.Build(
			beforeBreakdown, // before entity
			breakdownEntity, // after entity
			nil,
			"",
		)

		err = breakdownService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_BREAKDOWN,
				auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = breakdownService.FindAllPagination(&paginationRequest)
	return paginatedResp
}

func (breakdownService *ServiceImpl) Delete(ginContext *gin.Context, deleteBreakdownRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.BreakdownResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginatedResp *model.PaginatedResponse[*model.BreakdownResponse]
	valErr := breakdownService.validatorService.ValidateStruct(deleteBreakdownRequest)
	breakdownService.validatorService.ParseValidationError(valErr, *deleteBreakdownRequest)
	err := breakdownService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		breakdownEntity, err := breakdownService.breakdownRepository.FindById(gormTransaction, deleteBreakdownRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		breakdownEntity.Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		err = breakdownService.breakdownRepository.Delete(gormTransaction, deleteBreakdownRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := breakdownService.auditLogService.Build(
			breakdownEntity, // before entity
			nil,             // after entity
			nil,
			"",
		)

		err = breakdownService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_BREAKDOWN,
				auditPayload)
		return nil
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = breakdownService.FindAllPagination(&paginationRequest)
	return paginatedResp
}
