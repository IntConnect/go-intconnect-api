package machine

import (
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
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
) *ServiceImpl {
	return &ServiceImpl{
		machineRepository:         machineRepository,
		validatorService:          validatorService,
		dbConnection:              dbConnection,
		viperConfig:               viperConfig,
		localStorageService:       localStorageService,
		machineDocumentRepository: machineDocumentRepository,
		auditLogService:           auditLogService,
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

// Create - Membuat machine baru
func (machineService *ServiceImpl) Create(ginContext *gin.Context, createMachineRequest *model.CreateMachineRequest) *model.PaginatedResponse[*model.MachineResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	ipAddress, _ := helper.ExtractRequestMeta(ginContext)
	valErr := machineService.validatorService.ValidateStruct(createMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *createMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modelPath, err := machineService.localStorageService.Disk().Put(createMachineRequest.Model, fmt.Sprintf("machines/models/%d-%s", time.Now().UnixNano(), createMachineRequest.Model.Filename))
		fmt.Println(err)
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
		fmt.Println(err)
		thumbnailPath, err := machineService.localStorageService.Disk().Put(createMachineRequest.Thumbnail, fmt.Sprintf("machines/thumbnails/%d-%s", time.Now().UnixNano(), createMachineRequest.Thumbnail.Filename))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
		machineEntity := helper.MapCreateRequestIntoEntity[model.CreateMachineRequest, entity.Machine](createMachineRequest)
		machineEntity.ModelPath = modelPath
		machineEntity.ThumbnailPath = thumbnailPath
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
		machineService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      userJwtClaims.Id,
			Action:      model.AUDIT_LOG_CREATE,
			Feature:     model.AUDIT_LOG_FEATURE_MACHINE,
			Description: "",
			Before:      nil,
			After:       machineEntity,
			IpAddress:   ipAddress,
		})
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	var paginationResp *model.PaginatedResponse[*model.MachineResponse]
	paginationRequest := model.NewPaginationRequest()
	paginationResp = machineService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (machineService *ServiceImpl) Update(ginContext *gin.Context, updateMachineRequest *model.UpdateMachineRequest) *model.PaginatedResponse[*model.MachineResponse] {
	valErr := machineService.validatorService.ValidateStruct(updateMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *updateMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machine, err := machineService.machineRepository.FindById(gormTransaction, updateMachineRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateMachineRequest, machine)
		err = machineService.machineRepository.Update(gormTransaction, machine)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	var paginationResp *model.PaginatedResponse[*model.MachineResponse]
	paginationRequest := model.NewPaginationRequest()
	paginationResp = machineService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (machineService *ServiceImpl) Delete(ginContext *gin.Context, deleteMachineRequest *model.DeleteMachineRequest) *model.PaginatedResponse[*model.MachineResponse] {
	valErr := machineService.validatorService.ValidateStruct(deleteMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *deleteMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := machineService.machineRepository.Delete(gormTransaction, deleteMachineRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	var paginationResp *model.PaginatedResponse[*model.MachineResponse]
	paginationRequest := model.NewPaginationRequest()
	paginationResp = machineService.FindAllPagination(&paginationRequest)
	return paginationResp
}
