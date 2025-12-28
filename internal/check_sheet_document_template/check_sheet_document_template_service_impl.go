package check_sheet_document_template

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/parameter"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	checkSheetDocumentTemplateRepository Repository
	auditLogService                      auditLog.Service
	parameterRepository                  parameter.Repository
	validatorService                     validator.Service
	dbConnection                         *gorm.DB
	viperConfig                          *viper.Viper
}

func NewService(checkSheetDocumentTemplateRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	parameterRepository parameter.Repository,
	auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		checkSheetDocumentTemplateRepository: checkSheetDocumentTemplateRepository,
		validatorService:                     validatorService,
		dbConnection:                         dbConnection,
		viperConfig:                          viperConfig,
		parameterRepository:                  parameterRepository,
		auditLogService:                      auditLogService,
	}
}

func (checkSheetDocumentTemplateService *ServiceImpl) FindAll() []*model.CheckSheetDocumentTemplateResponse {
	var allCheckSheetDocumentTemplate []*model.CheckSheetDocumentTemplateResponse
	err := checkSheetDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetDocumentTemplateResponse, err := checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allCheckSheetDocumentTemplate = helper.MapEntitiesIntoResponsesWithFunc[*entity.CheckSheetDocumentTemplate, *model.CheckSheetDocumentTemplateResponse](checkSheetDocumentTemplateResponse, mapper.FuncMapAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allCheckSheetDocumentTemplate
}

func (checkSheetDocumentTemplateService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var checkSheetDocumentTemplateResponses []*model.CheckSheetDocumentTemplateResponse
	var totalItems int64

	err := checkSheetDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetDocumentTemplateEntities, total, err := checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		checkSheetDocumentTemplateResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.CheckSheetDocumentTemplate, *model.CheckSheetDocumentTemplateResponse](
			checkSheetDocumentTemplateEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"CheckSheet document templates fetched successfully",
		checkSheetDocumentTemplateResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat checkSheetDocumentTemplate baru
func (checkSheetDocumentTemplateService *ServiceImpl) Create(ginContext *gin.Context, createCheckSheetDocumentTemplateRequest *model.CreateCheckSheetDocumentTemplateRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetDocumentTemplateService.validatorService.ValidateStruct(createCheckSheetDocumentTemplateRequest)
	checkSheetDocumentTemplateService.validatorService.ParseValidationError(valErr, *createCheckSheetDocumentTemplateRequest)
	err := checkSheetDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetDocumentTemplateEntity := helper.MapCreateRequestIntoEntity[model.CreateCheckSheetDocumentTemplateRequest, entity.CheckSheetDocumentTemplate](createCheckSheetDocumentTemplateRequest)
		checkSheetDocumentTemplateEntity.Auditable = entity.NewAuditable("Administrator")
		parameterEntities, err := checkSheetDocumentTemplateService.parameterRepository.FindBatchById(gormTransaction, createCheckSheetDocumentTemplateRequest.ParameterIds)
		if len(parameterEntities) != len(createCheckSheetDocumentTemplateRequest.ParameterIds) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, exception.ErrSomeResourceNotFound))
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		checkSheetDocumentTemplateEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		err = checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.Create(gormTransaction, checkSheetDocumentTemplateEntity)
		var checkSheetDocumentTemplateParameters []*entity.CheckSheetDocumentTemplateParameter
		for _, parameterEntity := range parameterEntities {
			checkSheetDocumentTemplateParameters = append(checkSheetDocumentTemplateParameters, &entity.CheckSheetDocumentTemplateParameter{
				CheckSheetDocumentTemplateId: checkSheetDocumentTemplateEntity.Id,
				ParameterId:                  parameterEntity.Id,
			})
		}

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetDocumentTemplateService.auditLogService.Build(
			nil,                              // before entity
			checkSheetDocumentTemplateEntity, // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": nil,
					"after":  createCheckSheetDocumentTemplateRequest.ParameterIds,
				},
			},
			"",
		)

		err = checkSheetDocumentTemplateService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_CHECK_SHEET_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp := checkSheetDocumentTemplateService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (checkSheetDocumentTemplateService *ServiceImpl) Update(ginContext *gin.Context, updateCheckSheetDocumentTemplateRequest *model.UpdateCheckSheetDocumentTemplateRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetDocumentTemplateService.validatorService.ValidateStruct(updateCheckSheetDocumentTemplateRequest)
	checkSheetDocumentTemplateService.validatorService.ParseValidationError(valErr, *updateCheckSheetDocumentTemplateRequest)

	err := checkSheetDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetDocumentTemplate, err := checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.FindById(gormTransaction, updateCheckSheetDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		beforeCheckSheetDocumentTemplate := *checkSheetDocumentTemplate
		beforeCheckSheetDocumentTemplate.CheckSheetDocumentTemplateParameters = append([]*entity.CheckSheetDocumentTemplateParameter(nil), checkSheetDocumentTemplate.CheckSheetDocumentTemplateParameters...)
		checkSheetDocumentTemplate.Auditable = entity.NewAuditable(userJwtClaims.Username)

		// Map field biasa
		helper.MapUpdateRequestIntoEntity(updateCheckSheetDocumentTemplateRequest, checkSheetDocumentTemplate)
		if err = checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.Update(gormTransaction, checkSheetDocumentTemplate); err != nil {
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}

		// Generate daftar parameter baru
		newParameterIds := make([]entity.Parameter, 0, len(updateCheckSheetDocumentTemplateRequest.ParameterIds))
		for _, parameterId := range updateCheckSheetDocumentTemplateRequest.ParameterIds {
			newParameterIds = append(newParameterIds, entity.Parameter{Id: parameterId})
		}

		// Replace relasi M2M
		if err := gormTransaction.Model(checkSheetDocumentTemplate).Association("Parameters").Replace(&newParameterIds); err != nil {
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		auditPayload := checkSheetDocumentTemplateService.auditLogService.Build(
			&beforeCheckSheetDocumentTemplate, // before entity
			checkSheetDocumentTemplate,        // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": helper.ExtractIds(beforeCheckSheetDocumentTemplate.CheckSheetDocumentTemplateParameters),
					"after":  updateCheckSheetDocumentTemplateRequest.ParameterIds,
				},
			},
			"",
		)
		err = checkSheetDocumentTemplateService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_UPDATE,
				model.AUDIT_LOG_FEATURE_CHECK_SHEET_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationReq := model.NewPaginationRequest()
	return checkSheetDocumentTemplateService.FindAllPagination(&paginationReq)
}

func (checkSheetDocumentTemplateService *ServiceImpl) Delete(ginContext *gin.Context, deleteCheckSheetDocumentTemplateRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetDocumentTemplateService.validatorService.ValidateStruct(deleteCheckSheetDocumentTemplateRequest)
	checkSheetDocumentTemplateService.validatorService.ParseValidationError(valErr, *deleteCheckSheetDocumentTemplateRequest)
	err := checkSheetDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetDocumentTemplate, err := checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.FindById(gormTransaction, deleteCheckSheetDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		checkSheetDocumentTemplate.Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		err = checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.Delete(gormTransaction, deleteCheckSheetDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetDocumentTemplateService.auditLogService.Build(
			checkSheetDocumentTemplate, // before entity
			nil,                        // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": helper.ExtractIds(checkSheetDocumentTemplate.CheckSheetDocumentTemplateParameters),
					"after":  nil,
				},
			},
			deleteCheckSheetDocumentTemplateRequest.Reason,
		)
		err = checkSheetDocumentTemplateService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_DELETE,
				model.AUDIT_LOG_FEATURE_CHECK_SHEET_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationReq := model.NewPaginationRequest()
	return checkSheetDocumentTemplateService.FindAllPagination(&paginationReq)
}
