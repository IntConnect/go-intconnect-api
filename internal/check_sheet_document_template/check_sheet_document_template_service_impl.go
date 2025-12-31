package check_sheet_document_template

import (
	"fmt"
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
	checkSheetDocumentTemplateRepository Repository
	auditLogService                      auditLog.Service
	validatorService                     validator.Service
	dbConnection                         *gorm.DB
	viperConfig                          *viper.Viper
}

func NewService(checkSheetDocumentTemplateRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		checkSheetDocumentTemplateRepository: checkSheetDocumentTemplateRepository,
		validatorService:                     validatorService,
		dbConnection:                         dbConnection,
		viperConfig:                          viperConfig,
		auditLogService:                      auditLogService,
	}
}

func (checkSheetDocumentTemplateService *ServiceImpl) FindAll() []*model.CheckSheetDocumentTemplateResponse {
	var allCheckSheetDocumentTemplate []*model.CheckSheetDocumentTemplateResponse
	err := checkSheetDocumentTemplateService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetDocumentTemplateResponse, err := checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allCheckSheetDocumentTemplate = helper.MapEntitiesIntoResponsesWithFunc[*entity.CheckSheetDocumentTemplate, *model.CheckSheetDocumentTemplateResponse](checkSheetDocumentTemplateResponse, mapper.FuncMapAuditable, mapper.MapCheckSheetDocumentTemplate)
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
			mapper.MapCheckSheetDocumentTemplate,
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
		checkSheetDocumentTemplateEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		err := checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.Create(gormTransaction, checkSheetDocumentTemplateEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetDocumentTemplateService.auditLogService.Build(
			nil,
			checkSheetDocumentTemplateEntity,
			nil,
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
		checkSheetDocumentTemplate.Auditable = entity.UpdateAuditable(userJwtClaims.Username)
		checkSheetDocumentTemplate.RevisionNumber += 1
		// Map field biasa
		helper.MapUpdateRequestIntoEntity(updateCheckSheetDocumentTemplateRequest, checkSheetDocumentTemplate)
		fmt.Println(checkSheetDocumentTemplate)
		err = checkSheetDocumentTemplateService.checkSheetDocumentTemplateRepository.Update(gormTransaction, checkSheetDocumentTemplate)

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetDocumentTemplateService.auditLogService.Build(
			&beforeCheckSheetDocumentTemplate, // before entity
			checkSheetDocumentTemplate,        // after entity
			nil,
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
			nil,
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
