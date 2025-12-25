package check_sheet

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/parameter"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	checkSheetRepository Repository
	auditLogService      auditLog.Service
	parameterRepository  parameter.Repository
	validatorService     validator.Service
	dbConnection         *gorm.DB
	viperConfig          *viper.Viper
}

func NewService(checkSheetRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	parameterRepository parameter.Repository,
	auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		checkSheetRepository: checkSheetRepository,
		validatorService:     validatorService,
		dbConnection:         dbConnection,
		viperConfig:          viperConfig,
		parameterRepository:  parameterRepository,
		auditLogService:      auditLogService,
	}
}

func (checkSheetService *ServiceImpl) FindAll() []*model.CheckSheetResponse {
	var allCheckSheet []*model.CheckSheetResponse
	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetResponse, err := checkSheetService.checkSheetRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allCheckSheet = helper.MapEntitiesIntoResponsesWithFunc[*entity.CheckSheet, *model.CheckSheetResponse](checkSheetResponse, mapper.FuncMapAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allCheckSheet
}

func (checkSheetService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.CheckSheetResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var checkSheetResponses []*model.CheckSheetResponse
	var totalItems int64

	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetEntities, total, err := checkSheetService.checkSheetRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		checkSheetResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.CheckSheet, *model.CheckSheetResponse](
			checkSheetEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"CheckSheet document templates fetched successfully",
		checkSheetResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat checkSheet baru
func (checkSheetService *ServiceImpl) Create(ginContext *gin.Context, createCheckSheetRequest *model.CreateCheckSheetRequest) *model.PaginatedResponse[*model.CheckSheetResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetService.validatorService.ValidateStruct(createCheckSheetRequest)
	checkSheetService.validatorService.ParseValidationError(valErr, *createCheckSheetRequest)
	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetEntity := helper.MapCreateRequestIntoEntity[model.CreateCheckSheetRequest, entity.CheckSheet](createCheckSheetRequest)
		checkSheetEntity.Auditable = entity.NewAuditable("Administrator")
		err := checkSheetService.checkSheetRepository.Create(gormTransaction, checkSheetEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetService.auditLogService.Build(
			nil,              // before entity
			checkSheetEntity, // after entity
			nil,
			"",
		)

		err = checkSheetService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_CHECK_SHEET_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp := checkSheetService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (checkSheetService *ServiceImpl) Update(ginContext *gin.Context, updateCheckSheetRequest *model.UpdateCheckSheetRequest) *model.PaginatedResponse[*model.CheckSheetResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetService.validatorService.ValidateStruct(updateCheckSheetRequest)
	checkSheetService.validatorService.ParseValidationError(valErr, *updateCheckSheetRequest)

	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheet, err := checkSheetService.checkSheetRepository.FindById(gormTransaction, updateCheckSheetRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		beforeCheckSheet := *checkSheet
		beforeCheckSheet.Parameters = append([]*entity.Parameter(nil), checkSheet.Parameters...)
		checkSheet.Auditable = entity.NewAuditable(userJwtClaims.Username)

		// Map field biasa
		helper.MapUpdateRequestIntoEntity(updateCheckSheetRequest, checkSheet)
		if err = checkSheetService.checkSheetRepository.Update(gormTransaction, checkSheet); err != nil {
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}

		// Generate daftar parameter baru
		newParameterIds := make([]entity.Parameter, 0, len(updateCheckSheetRequest.ParameterIds))
		for _, parameterId := range updateCheckSheetRequest.ParameterIds {
			newParameterIds = append(newParameterIds, entity.Parameter{Id: parameterId})
		}

		// Replace relasi M2M
		if err := gormTransaction.Model(checkSheet).Association("Parameters").Replace(&newParameterIds); err != nil {
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		auditPayload := checkSheetService.auditLogService.Build(
			&beforeCheckSheet, // before entity
			checkSheet,        // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": helper.ExtractIds(beforeCheckSheet.Parameters),
					"after":  updateCheckSheetRequest.ParameterIds,
				},
			},
			"",
		)
		err = checkSheetService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_UPDATE,
				model.AUDIT_LOG_FEATURE_CHECK_SHEET_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationReq := model.NewPaginationRequest()
	return checkSheetService.FindAllPagination(&paginationReq)
}

func (checkSheetService *ServiceImpl) Delete(ginContext *gin.Context, deleteCheckSheetRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.CheckSheetResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetService.validatorService.ValidateStruct(deleteCheckSheetRequest)
	checkSheetService.validatorService.ParseValidationError(valErr, *deleteCheckSheetRequest)
	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheet, err := checkSheetService.checkSheetRepository.FindById(gormTransaction, deleteCheckSheetRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		checkSheet.Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		err = checkSheetService.checkSheetRepository.Delete(gormTransaction, deleteCheckSheetRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetService.auditLogService.Build(
			checkSheet, // before entity
			nil,        // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": helper.ExtractIds(checkSheet.Parameters),
					"after":  nil,
				},
			},
			deleteCheckSheetRequest.Reason,
		)
		err = checkSheetService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_DELETE,
				model.AUDIT_LOG_FEATURE_CHECK_SHEET_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationReq := model.NewPaginationRequest()
	return checkSheetService.FindAllPagination(&paginationReq)
}
