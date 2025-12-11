package report_document_template

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
	reportDocumentTemplateRepository Repository
	auditLogService                  auditLog.Service
	parameterRepository              parameter.Repository
	validatorService                 validator.Service
	dbConnection                     *gorm.DB
	viperConfig                      *viper.Viper
}

func NewService(reportDocumentTemplateRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	parameterRepository parameter.Repository,
	auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		reportDocumentTemplateRepository: reportDocumentTemplateRepository,
		validatorService:                 validatorService,
		dbConnection:                     dbConnection,
		viperConfig:                      viperConfig,
		parameterRepository:              parameterRepository,
		auditLogService:                  auditLogService,
	}
}

func (reportDocumentTemplateService *ServiceImpl) FindAll() []*model.ReportDocumentTemplateResponse {
	var allReportDocumentTemplate []*model.ReportDocumentTemplateResponse
	err := reportDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		reportDocumentTemplateResponse, err := reportDocumentTemplateService.reportDocumentTemplateRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allReportDocumentTemplate = helper.MapEntitiesIntoResponsesWithFunc[entity.ReportDocumentTemplate, *model.ReportDocumentTemplateResponse](reportDocumentTemplateResponse, mapper.FuncMapAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allReportDocumentTemplate
}

func (reportDocumentTemplateService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var reportDocumentTemplateResponses []*model.ReportDocumentTemplateResponse
	var totalItems int64

	err := reportDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		reportDocumentTemplateEntities, total, err := reportDocumentTemplateService.reportDocumentTemplateRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		reportDocumentTemplateResponses = helper.MapEntitiesIntoResponsesWithFunc[entity.ReportDocumentTemplate, *model.ReportDocumentTemplateResponse](
			reportDocumentTemplateEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Report document templates fetched successfully",
		reportDocumentTemplateResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat reportDocumentTemplate baru
func (reportDocumentTemplateService *ServiceImpl) Create(ginContext *gin.Context, createReportDocumentTemplateRequest *model.CreateReportDocumentTemplateRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := reportDocumentTemplateService.validatorService.ValidateStruct(createReportDocumentTemplateRequest)
	reportDocumentTemplateService.validatorService.ParseValidationError(valErr, *createReportDocumentTemplateRequest)
	err := reportDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		reportDocumentTemplateEntity := helper.MapCreateRequestIntoEntity[model.CreateReportDocumentTemplateRequest, entity.ReportDocumentTemplate](createReportDocumentTemplateRequest)
		reportDocumentTemplateEntity.Auditable = entity.NewAuditable("Administrator")
		parameterEntities, err := reportDocumentTemplateService.parameterRepository.FindBatchById(gormTransaction, createReportDocumentTemplateRequest.ParameterIds)
		if len(parameterEntities) != len(createReportDocumentTemplateRequest.ParameterIds) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, exception.ErrSomeResourceNotFound))
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		reportDocumentTemplateEntity.Parameters = parameterEntities
		reportDocumentTemplateEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		err = reportDocumentTemplateService.reportDocumentTemplateRepository.Create(gormTransaction, reportDocumentTemplateEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := reportDocumentTemplateService.auditLogService.Build(
			nil,                          // before entity
			reportDocumentTemplateEntity, // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": nil,
					"after":  createReportDocumentTemplateRequest.ParameterIds,
				},
			},
			"",
		)

		err = reportDocumentTemplateService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_REPORT_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp := reportDocumentTemplateService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (reportDocumentTemplateService *ServiceImpl) Update(ginContext *gin.Context, updateReportDocumentTemplateRequest *model.UpdateReportDocumentTemplateRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := reportDocumentTemplateService.validatorService.ValidateStruct(updateReportDocumentTemplateRequest)
	reportDocumentTemplateService.validatorService.ParseValidationError(valErr, *updateReportDocumentTemplateRequest)

	err := reportDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		reportDocumentTemplate, err := reportDocumentTemplateService.reportDocumentTemplateRepository.FindById(gormTransaction, updateReportDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		beforeReportDocumentTemplate := *reportDocumentTemplate
		beforeReportDocumentTemplate.Parameters = append([]*entity.Parameter(nil), reportDocumentTemplate.Parameters...)
		reportDocumentTemplate.Auditable = entity.NewAuditable(userJwtClaims.Username)

		// Map field biasa
		helper.MapUpdateRequestIntoEntity(updateReportDocumentTemplateRequest, reportDocumentTemplate)
		if err = reportDocumentTemplateService.reportDocumentTemplateRepository.Update(gormTransaction, reportDocumentTemplate); err != nil {
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}

		// Generate daftar parameter baru
		newParameterIds := make([]entity.Parameter, 0, len(updateReportDocumentTemplateRequest.ParameterIds))
		for _, parameterId := range updateReportDocumentTemplateRequest.ParameterIds {
			newParameterIds = append(newParameterIds, entity.Parameter{Id: parameterId})
		}

		// Replace relasi M2M
		if err := gormTransaction.Model(reportDocumentTemplate).Association("Parameters").Replace(&newParameterIds); err != nil {
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		auditPayload := reportDocumentTemplateService.auditLogService.Build(
			&beforeReportDocumentTemplate, // before entity
			reportDocumentTemplate,        // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": helper.ExtractIds(beforeReportDocumentTemplate.Parameters),
					"after":  updateReportDocumentTemplateRequest.ParameterIds,
				},
			},
			"",
		)
		err = reportDocumentTemplateService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_UPDATE,
				model.AUDIT_LOG_FEATURE_REPORT_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationReq := model.NewPaginationRequest()
	return reportDocumentTemplateService.FindAllPagination(&paginationReq)
}

func (reportDocumentTemplateService *ServiceImpl) Delete(ginContext *gin.Context, deleteReportDocumentTemplateRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := reportDocumentTemplateService.validatorService.ValidateStruct(deleteReportDocumentTemplateRequest)
	reportDocumentTemplateService.validatorService.ParseValidationError(valErr, *deleteReportDocumentTemplateRequest)
	err := reportDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		reportDocumentTemplate, err := reportDocumentTemplateService.reportDocumentTemplateRepository.FindById(gormTransaction, deleteReportDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		reportDocumentTemplate.Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		err = reportDocumentTemplateService.reportDocumentTemplateRepository.Delete(gormTransaction, deleteReportDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := reportDocumentTemplateService.auditLogService.Build(
			reportDocumentTemplate, // before entity
			nil,                    // after entity
			map[string]map[string][]uint64{
				"parameters": {
					"before": helper.ExtractIds(reportDocumentTemplate.Parameters),
					"after":  nil,
				},
			},
			deleteReportDocumentTemplateRequest.Reason,
		)
		err = reportDocumentTemplateService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_DELETE,
				model.AUDIT_LOG_FEATURE_REPORT_DOCUMENT_TEMPLATE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationReq := model.NewPaginationRequest()
	return reportDocumentTemplateService.FindAllPagination(&paginationReq)
}
