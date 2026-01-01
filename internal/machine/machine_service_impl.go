package machine

import (
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
	dashboardWidget "go-intconnect-api/internal/dashboard_widget"
	"go-intconnect-api/internal/entity"
	machineDocument "go-intconnect-api/internal/machine_document"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/storage"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	machineRepository         Repository
	dashboardWidgetRepository dashboardWidget.Repository
	auditLogService           auditLog.Service
	machineDocumentRepository machineDocument.Repository
	validatorService          validator.Service
	dbConnection              *gorm.DB
	viperConfig               *viper.Viper
	localStorageService       *storage.Manager
}

func NewService(machineRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	localStorageService *storage.Manager,
	machineDocumentRepository machineDocument.Repository,
	auditLogService auditLog.Service,
	dashboardWidgetRepository dashboardWidget.Repository,
) *ServiceImpl {
	return &ServiceImpl{
		machineRepository:         machineRepository,
		validatorService:          validatorService,
		dbConnection:              dbConnection,
		viperConfig:               viperConfig,
		localStorageService:       localStorageService,
		machineDocumentRepository: machineDocumentRepository,
		auditLogService:           auditLogService,
		dashboardWidgetRepository: dashboardWidgetRepository,
	}
}

func (machineService *ServiceImpl) FindAll() []*model.MachineResponse {
	var machineResponsesRequest []*model.MachineResponse
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntities, err := machineService.machineRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		machineResponsesRequest = helper.MapEntitiesIntoResponsesWithFunc[entity.Machine, *model.MachineResponse](machineEntities, mapper.FuncMapAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return machineResponsesRequest
}

func (machineService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MachineResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var machineResponses []*model.MachineResponse
	var totalItems int64

	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntities, total, err := machineService.machineRepository.FindAllPagination(gormTransaction, paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		machineResponses = helper.MapEntitiesIntoResponsesWithFunc[entity.Machine, *model.MachineResponse](
			machineEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Machines fetched successfully",
		machineResponses,
		paginationReq,
		totalItems,
	)
}

func (machineService *ServiceImpl) FindById(ginContext *gin.Context, machineId uint64) *model.MachineResponse {
	var machineResponseRequest *model.MachineResponse
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntity, err := machineService.machineRepository.FindById(gormTransaction, machineId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		machineResponseRequest = helper.MapEntityIntoResponse[*entity.Machine, *model.MachineResponse](machineEntity,
			mapper.FuncMapAuditable, mapper.MapMachineDocument)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return machineResponseRequest
}

func (machineService *ServiceImpl) FindByFacilityId(ginContext *gin.Context, facilityId uint64) []*model.MachineResponse {
	var machineResponseRequest []*model.MachineResponse
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntities, err := machineService.machineRepository.FindByFacilityId(gormTransaction, facilityId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		machineResponseRequest = helper.MapEntitiesIntoResponsesWithFunc[*entity.Machine, *model.MachineResponse](machineEntities, mapper.FuncMapAuditable, mapper.MapMachineDocument)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return machineResponseRequest
}

func (machineService *ServiceImpl) ManageDashboard(ginContext *gin.Context, machineDashboardWidget *model.MachineDashboardWidget) {
	valErr := machineService.validatorService.ValidateStruct(machineDashboardWidget)
	machineService.validatorService.ParseValidationError(valErr, *machineDashboardWidget)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var dashboardWidgetEntities []*entity.DashboardWidget
		for _, dashboardWidgetRequest := range machineDashboardWidget.DashboardWidget {
			var dashboardWidgetEntity entity.DashboardWidget
			helper.MapUpdateRequestIntoEntity(dashboardWidgetRequest, &dashboardWidgetEntity)
			dashboardWidgetEntity.MachineId = machineDashboardWidget.MachineId
			dashboardWidgetEntities = append(dashboardWidgetEntities, &dashboardWidgetEntity)
		}
		if len(dashboardWidgetEntities) > 0 {
			err := machineService.dashboardWidgetRepository.CreateBatch(gormTransaction, dashboardWidgetEntities)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

// Create - Membuat machine baru
func (machineService *ServiceImpl) Create(ginContext *gin.Context, createMachineRequest *model.CreateMachineRequest) {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := machineService.validatorService.ValidateStruct(createMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *createMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modelPath, err := machineService.localStorageService.Disk().Put(createMachineRequest.Model, fmt.Sprintf("machines/models/%d-%s", time.Now().UnixNano(), createMachineRequest.Model.Filename))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
		thumbnailPath, err := machineService.localStorageService.Disk().Put(createMachineRequest.Thumbnail, fmt.Sprintf("machines/thumbnails/%d-%s", time.Now().UnixNano(), createMachineRequest.Thumbnail.Filename))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
		machineEntity := helper.MapCreateRequestIntoEntity[model.CreateMachineRequest, entity.Machine](createMachineRequest)
		machineEntity.ModelPath = modelPath
		machineEntity.ThumbnailPath = thumbnailPath
		machineEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		err = machineService.machineRepository.Create(gormTransaction, machineEntity)
		var machineDocumentEntities []*entity.MachineDocument
		for _, createMachineDocumentRequest := range createMachineRequest.MachineDocuments {
			machineDocumentEntity := helper.MapCreateRequestIntoEntity[model.CreateMachineDocumentRequest, entity.MachineDocument](&createMachineDocumentRequest)
			machineDocumentFilePath, err := machineService.localStorageService.Disk().Put(createMachineDocumentRequest.DocumentFile, fmt.Sprintf("machines/documents/%d-%s", time.Now().UnixNano(), createMachineDocumentRequest.DocumentFile.Filename))
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
			machineDocumentEntity.FilePath = machineDocumentFilePath
			machineDocumentEntity.MachineId = machineEntity.Id
			machineDocumentEntities = append(machineDocumentEntities, machineDocumentEntity)
		}
		if len(machineDocumentEntities) > 0 {
			err = machineService.machineDocumentRepository.CreateBatch(gormTransaction, machineDocumentEntities)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		auditPayload := machineService.auditLogService.Build(
			nil,           // before entity
			machineEntity, // after entity
			map[string]map[string][]uint64{
				"machine_documents": {
					"before": nil,
					"after":  helper.ExtractIds(machineDocumentEntities),
				},
			},
			"",
		)

		err = machineService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_MACHINE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))

}

func (machineService *ServiceImpl) Update(ginContext *gin.Context, updateMachineRequest *model.UpdateMachineRequest) {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := machineService.validatorService.ValidateStruct(updateMachineRequest)

	machineService.validatorService.ParseValidationError(valErr, *updateMachineRequest)

	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntity := helper.MapCreateRequestIntoEntity[model.UpdateMachineRequest, entity.Machine](updateMachineRequest)
		if updateMachineRequest.Model != nil {
			modelPath, err := machineService.localStorageService.Disk().Put(updateMachineRequest.Model, fmt.Sprintf("machines/models/%d-%s", time.Now().UnixNano(), updateMachineRequest.Model.Filename))
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
			machineEntity.ModelPath = modelPath
		}
		if updateMachineRequest.Thumbnail != nil {
			thumbnailPath, err := machineService.localStorageService.Disk().Put(updateMachineRequest.Thumbnail, fmt.Sprintf("machines/thumbnails/%d-%s", time.Now().UnixNano(), updateMachineRequest.Thumbnail.Filename))
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
			machineEntity.ThumbnailPath = thumbnailPath
		}
		machineEntity.Auditable = entity.UpdateAuditable(userJwtClaims.Username)
		err := machineService.machineRepository.Update(gormTransaction, machineEntity)
		machineDocuments, err := machineService.machineDocumentRepository.FindBatchById(gormTransaction, updateMachineRequest.DeletedMachineDocumentIds)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		if len(machineDocuments) != len(updateMachineRequest.DeletedMachineDocumentIds) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusNotFound, fmt.Sprintf("%s machine documents", exception.ErrSomeResourceNotFound)))
		}
		for i, machineDocumentEntity := range machineDocuments {
			newPath, err := machineService.localStorageService.Disk().MoveFile(machineDocumentEntity.FilePath, "machines/documents")
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrSavingResources))
			machineDocuments[i].FilePath = newPath
			machineDocuments[i].Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		}
		var machineDocumentEntities []*entity.MachineDocument
		for _, createMachineDocumentRequest := range updateMachineRequest.MachineDocuments {
			machineDocumentEntity := helper.MapCreateRequestIntoEntity[model.CreateMachineDocumentRequest, entity.MachineDocument](&createMachineDocumentRequest)
			machineDocumentFilePath, err := machineService.localStorageService.Disk().Put(createMachineDocumentRequest.DocumentFile, fmt.Sprintf("machines/documents/%d-%s", time.Now().UnixNano(), createMachineDocumentRequest.DocumentFile.Filename))
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
			machineDocumentEntity.FilePath = machineDocumentFilePath
			machineDocumentEntity.MachineId = machineEntity.Id
			machineDocumentEntities = append(machineDocumentEntities, machineDocumentEntity)
		}
		if len(machineDocumentEntities) > 0 {
			err = machineService.machineDocumentRepository.CreateBatch(gormTransaction, machineDocumentEntities)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		auditPayload := machineService.auditLogService.Build(
			nil,           // before entity
			machineEntity, // after entity
			map[string]map[string][]uint64{
				"machine_documents": {
					"before": helper.ExtractIds(machineDocuments),
					"after":  helper.ExtractIds(machineDocumentEntities),
				},
			},
			"",
		)

		err = machineService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_UPDATE,
				model.AUDIT_LOG_FEATURE_MACHINE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))

}

func (machineService *ServiceImpl) Delete(ginContext *gin.Context, deleteMachineRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.MachineResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := machineService.validatorService.ValidateStruct(deleteMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *deleteMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntity, err := machineService.machineRepository.FindById(gormTransaction, deleteMachineRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		machineEntity.Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		for i, machineDocumentEntity := range machineEntity.MachineDocuments {
			newPath, err := machineService.localStorageService.Disk().MoveFile(machineDocumentEntity.FilePath, "machines/documents")
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrSavingResources))
			machineEntity.MachineDocuments[i].FilePath = newPath
			machineEntity.MachineDocuments[i].Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		}
		newPath, err := machineService.localStorageService.Disk().MoveFile(machineEntity.ThumbnailPath, "machines/thumbnails")
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrSavingResources))
		machineEntity.ThumbnailPath = newPath
		newPath, err = machineService.localStorageService.Disk().MoveFile(machineEntity.ModelPath, "machines/models")
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrSavingResources))
		machineEntity.ModelPath = newPath
		err = machineService.machineRepository.Delete(gormTransaction, machineEntity)

		auditPayload := machineService.auditLogService.Build(
			machineEntity, // before entity
			nil,           // after entity
			nil,
			deleteMachineRequest.Reason,
		)

		err = machineService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_DELETE,
				model.AUDIT_LOG_FEATURE_MACHINE,
				auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	var paginationResp *model.PaginatedResponse[*model.MachineResponse]
	paginationRequest := model.NewPaginationRequest()
	paginationResp = machineService.FindAllPagination(&paginationRequest)
	return paginationResp
}
