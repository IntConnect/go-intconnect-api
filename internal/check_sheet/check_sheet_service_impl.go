package check_sheet

import (
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
	checkSheetCheckPoint "go-intconnect-api/internal/check_sheet_check_point"
	checkSheetCheckPointValue "go-intconnect-api/internal/check_sheet_check_point_value"
	checkSheetDocumentTemplate "go-intconnect-api/internal/check_sheet_document_template"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/parameter"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	checkSheetRepository                 Repository
	checkSheetDocumentTemplateRepository checkSheetDocumentTemplate.Repository
	checkSheetCheckPointRepository       checkSheetCheckPoint.Repository
	checkSheetCheckPointValueRepository  checkSheetCheckPointValue.Repository
	auditLogService                      auditLog.Service
	parameterRepository                  parameter.Repository
	validatorService                     validator.Service
	dbConnection                         *gorm.DB
	viperConfig                          *viper.Viper
}

func NewService(checkSheetRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	parameterRepository parameter.Repository,
	auditLogService auditLog.Service,
	checkSheetDocumentTemplateRepository checkSheetDocumentTemplate.Repository,
	checkSheetCheckPointRepository checkSheetCheckPoint.Repository,
	checkSheetCheckPointValueRepository checkSheetCheckPointValue.Repository,

) *ServiceImpl {
	return &ServiceImpl{
		checkSheetRepository:                 checkSheetRepository,
		validatorService:                     validatorService,
		dbConnection:                         dbConnection,
		viperConfig:                          viperConfig,
		parameterRepository:                  parameterRepository,
		auditLogService:                      auditLogService,
		checkSheetDocumentTemplateRepository: checkSheetDocumentTemplateRepository,
		checkSheetCheckPointRepository:       checkSheetCheckPointRepository,
		checkSheetCheckPointValueRepository:  checkSheetCheckPointValueRepository,
	}

}

func (checkSheetService *ServiceImpl) FindAll() []*model.CheckSheetResponse {
	var allCheckSheet []*model.CheckSheetResponse
	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetResponse, err := checkSheetService.checkSheetRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allCheckSheet = helper.MapEntitiesIntoResponsesWithFunc[*entity.CheckSheet, *model.CheckSheetResponse](checkSheetResponse, mapper.FuncMapAuditable, mapper.MapCheckSheet)
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
			mapper.MapCheckSheet,
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

func (checkSheetService *ServiceImpl) FindById(ginContext *gin.Context, checkSheetId uint64) *model.CheckSheetResponse {
	var checkSheetResponse *model.CheckSheetResponse
	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetEntity, err := checkSheetService.checkSheetRepository.FindById(gormTransaction, checkSheetId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		checkSheetResponse = helper.MapEntityIntoResponse[*entity.CheckSheet, *model.CheckSheetResponse](checkSheetEntity, mapper.FuncMapAuditable)
		checkSheetResponse.CheckSheetCheckPoints = helper.MapEntitiesIntoResponses[*entity.CheckSheetCheckPoint, *model.CheckSheetCheckPointResponse](checkSheetEntity.CheckSheetCheckPoint)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return checkSheetResponse
}

func (checkSheetService *ServiceImpl) Create(ginContext *gin.Context, createCheckSheetRequest *model.CreateCheckSheetRequest) *model.PaginatedResponse[*model.CheckSheetResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetService.validatorService.ValidateStruct(createCheckSheetRequest)
	checkSheetService.validatorService.ParseValidationError(valErr, *createCheckSheetRequest)
	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetEntity := helper.MapCreateRequestIntoEntity[model.CreateCheckSheetRequest, entity.CheckSheet](createCheckSheetRequest)
		checkSheetEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		checkSheetEntity.ReportedBy = userJwtClaims.Id
		checkSheetEntity.Timestamp = time.Now()
		checkSheetEntity.Status = "Draft"
		checkSheetEntity.CheckSheetCheckPoint = nil
		err := checkSheetService.checkSheetRepository.Create(gormTransaction, checkSheetEntity)
		fmt.Print(err)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		var checkSheetCheckPointEntities []*entity.CheckSheetCheckPoint
		for _, checkSheetCheckPointRequest := range createCheckSheetRequest.CheckSheetCheckPoints {

			checkSheetCheckPointEntities = append(checkSheetCheckPointEntities, &entity.CheckSheetCheckPoint{
				CheckSheetId: checkSheetEntity.Id,
				ParameterId:  checkSheetCheckPointRequest.ParameterId,
				Name:         checkSheetCheckPointRequest.Name,
				Auditable:    entity.NewAuditable(userJwtClaims.Username),
			})
		}
		err = checkSheetService.checkSheetCheckPointRepository.CreateBatch(gormTransaction, checkSheetCheckPointEntities)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.DebugArrPointer(checkSheetCheckPointEntities)
		var checkSheetCheckPointValueEntities []*entity.CheckSheetCheckPointValue

		for i, cp := range checkSheetCheckPointEntities {
			for _, v := range createCheckSheetRequest.CheckSheetCheckPoints[i].CheckSheetValues {
				checkSheetCheckPointValueEntities = append(checkSheetCheckPointValueEntities, &entity.CheckSheetCheckPointValue{
					CheckSheetCheckPointId: cp.Id,
					Timestamp:              v.Timestamp,
					Value:                  v.Value,
				})
			}
		}

		if len(checkSheetCheckPointValueEntities) > 0 {
			err = checkSheetService.checkSheetCheckPointValueRepository.
				CreateBatch(gormTransaction, checkSheetCheckPointValueEntities)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}

		auditPayload := checkSheetService.auditLogService.Build(
			nil,              // before entity
			checkSheetEntity, // after entity
			map[string]map[string][]uint64{
				"check_sheet_checkpoint_id": {
					"before": nil,
					"after":  helper.ExtractIds(checkSheetCheckPointEntities),
				},
			},
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
func (checkSheetService *ServiceImpl) Approval(ginContext *gin.Context, approvalCheckSheet *model.ApprovalCheckSheet) *model.PaginatedResponse[*model.CheckSheetResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetService.validatorService.ValidateStruct(approvalCheckSheet)
	checkSheetService.validatorService.ParseValidationError(valErr, *approvalCheckSheet)
	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		decisionString := "Rejected"
		if approvalCheckSheet.Decision {
			decisionString = "Approved"
		}
		checkSheetEntity, err := checkSheetService.checkSheetRepository.FindById(gormTransaction, approvalCheckSheet.CheckSheetId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		checkSheetEntity.Status = decisionString
		checkSheetEntity.VerifiedBy = &userJwtClaims.Id
		checkSheetEntity.Auditable = entity.UpdateAuditable(userJwtClaims.Username)
		checkSheetEntity.Note = approvalCheckSheet.Note
		err = checkSheetService.checkSheetRepository.Update(gormTransaction, checkSheetEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetService.auditLogService.Build(
			nil,              // before entity
			checkSheetEntity, // after entity
			nil,
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
	paginationRequest := model.NewPaginationRequest()
	paginationResp := checkSheetService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (checkSheetService *ServiceImpl) Update(ginContext *gin.Context, updateCheckSheetRequest *model.UpdateCheckSheetRequest) *model.PaginatedResponse[*model.CheckSheetResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := checkSheetService.validatorService.ValidateStruct(updateCheckSheetRequest)
	checkSheetService.validatorService.ParseValidationError(valErr, *updateCheckSheetRequest)

	err := checkSheetService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		checkSheetEntity, err := checkSheetService.checkSheetRepository.FindById(gormTransaction, updateCheckSheetRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity[*model.UpdateCheckSheetRequest, entity.CheckSheet](updateCheckSheetRequest, checkSheetEntity)
		checkSheetEntity.Auditable = entity.UpdateAuditable(userJwtClaims.Username)
		err = checkSheetService.checkSheetRepository.Update(gormTransaction, checkSheetEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = checkSheetService.checkSheetCheckPointRepository.DeleteBatchById(gormTransaction, checkSheetEntity.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		var checkSheetCheckPointEntities []*entity.CheckSheetCheckPoint
		for _, checkSheetCheckPointRequest := range updateCheckSheetRequest.CheckSheetCheckPoint {
			var checkSheetCheckPointValueEntities []*entity.CheckSheetCheckPointValue
			for _, checkSheetValue := range checkSheetCheckPointRequest.CheckSheetValues {
				checkSheetCheckPointValueEntities = append(checkSheetCheckPointValueEntities, &entity.CheckSheetCheckPointValue{
					Timestamp: checkSheetValue.Timestamp,
					Value:     checkSheetValue.Value,
				})
			}
			checkSheetCheckPointEntities = append(checkSheetCheckPointEntities, &entity.CheckSheetCheckPoint{
				CheckSheetId:              checkSheetEntity.Id,
				ParameterId:               checkSheetCheckPointRequest.ParameterId,
				Name:                      checkSheetCheckPointRequest.Name,
				Auditable:                 entity.NewAuditable(userJwtClaims.Username),
				CheckSheetCheckPointValue: checkSheetCheckPointValueEntities,
			})
		}
		err = checkSheetService.checkSheetCheckPointRepository.CreateBatch(gormTransaction, checkSheetCheckPointEntities)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := checkSheetService.auditLogService.Build(
			nil,              // before entity
			checkSheetEntity, // after entity
			map[string]map[string][]uint64{
				"check_sheet_value_id": {
					"before": helper.ExtractIds(checkSheetEntity.CheckSheetCheckPoint),
					"after":  helper.ExtractIds(updateCheckSheetRequest.CheckSheetCheckPoint),
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
				"check_sheet_value_id": {
					"before": helper.ExtractIds(checkSheet.CheckSheetCheckPoint),
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
