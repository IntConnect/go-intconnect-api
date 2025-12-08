package report_document_template

import (
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
	parameterRepository              parameter.Repository
	validatorService                 validator.Service
	dbConnection                     *gorm.DB
	viperConfig                      *viper.Viper
}

func NewService(reportDocumentTemplateRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	parameterRepository parameter.Repository) *ServiceImpl {
	return &ServiceImpl{
		reportDocumentTemplateRepository: reportDocumentTemplateRepository,
		validatorService:                 validatorService,
		dbConnection:                     dbConnection,
		viperConfig:                      viperConfig,
		parameterRepository:              parameterRepository,
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
		err = reportDocumentTemplateService.reportDocumentTemplateRepository.Create(gormTransaction, reportDocumentTemplateEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp := reportDocumentTemplateService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (reportDocumentTemplateService *ServiceImpl) Update(ginContext *gin.Context, updateReportDocumentTemplateRequest *model.UpdateReportDocumentTemplateRequest) {
	valErr := reportDocumentTemplateService.validatorService.ValidateStruct(updateReportDocumentTemplateRequest)
	reportDocumentTemplateService.validatorService.ParseValidationError(valErr, *updateReportDocumentTemplateRequest)
	err := reportDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		reportDocumentTemplate, err := reportDocumentTemplateService.reportDocumentTemplateRepository.FindById(gormTransaction, updateReportDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateReportDocumentTemplateRequest, reportDocumentTemplate)
		err = reportDocumentTemplateService.reportDocumentTemplateRepository.Update(gormTransaction, reportDocumentTemplate)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (reportDocumentTemplateService *ServiceImpl) Delete(ginContext *gin.Context, deleteReportDocumentTemplateRequest *model.DeleteResourceGeneralRequest) {
	valErr := reportDocumentTemplateService.validatorService.ValidateStruct(deleteReportDocumentTemplateRequest)
	reportDocumentTemplateService.validatorService.ParseValidationError(valErr, *deleteReportDocumentTemplateRequest)
	err := reportDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := reportDocumentTemplateService.reportDocumentTemplateRepository.Delete(gormTransaction, deleteReportDocumentTemplateRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
