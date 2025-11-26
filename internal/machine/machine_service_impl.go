package machine

import (
	"fmt"
	"go-intconnect-api/internal/entity"
	machineDocument "go-intconnect-api/internal/machine_document"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/storage"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	machineRepository         Repository
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
) *ServiceImpl {
	return &ServiceImpl{
		machineRepository:         machineRepository,
		validatorService:          validatorService,
		dbConnection:              dbConnection,
		viperConfig:               viperConfig,
		localStorageService:       localStorageService,
		machineDocumentRepository: machineDocumentRepository,
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
func (machineService *ServiceImpl) Create(ginContext *gin.Context, createMachineRequest *model.CreateMachineRequest, modelFile *multipart.FileHeader, thumbnailFile *multipart.FileHeader) {
	valErr := machineService.validatorService.ValidateStruct(createMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *createMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		modelPath, err := machineService.localStorageService.Disk().Put(modelFile, fmt.Sprintf("machines/models/%d-%s", time.Now().UnixNano(), modelFile.Filename))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources, err))
		thumbnailPath, err := machineService.localStorageService.Disk().Put(thumbnailFile, fmt.Sprintf("machines/thumbnails/%d-%s", time.Now().UnixNano(), thumbnailFile.Filename))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources, err))
		machineEntity := helper.MapCreateRequestIntoEntity[model.CreateMachineRequest, entity.Machine](createMachineRequest)
		machineEntity.ModelPath = modelPath
		machineEntity.ThumbnailPath = thumbnailPath
		err = machineService.machineRepository.Create(gormTransaction, machineEntity)
		var machineDocumentEntities []*entity.MachineDocument
		fmt.Println(createMachineRequest.MachineDocuments)
		for _, createMachineDocumentRequest := range createMachineRequest.MachineDocuments {
			machineDocumentEntity := helper.MapCreateRequestIntoEntity[model.CreateMachineDocumentRequest, entity.MachineDocument](&createMachineDocumentRequest)
			machineDocumentFilePath, err := machineService.localStorageService.Disk().Put(createMachineDocumentRequest.DocumentFile, fmt.Sprintf("machines/documents/%d-%s", time.Now().UnixNano(), createMachineDocumentRequest.DocumentFile.Filename))
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources, err))
			machineDocumentEntity.FilePath = machineDocumentFilePath
			machineDocumentEntity.MachineId = machineEntity.Id
			machineDocumentEntities = append(machineDocumentEntities, machineDocumentEntity)
		}
		err = machineService.machineDocumentRepository.CreateBatch(gormTransaction, machineDocumentEntities)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (machineService *ServiceImpl) Update(ginContext *gin.Context, updateMachineRequest *model.UpdateMachineRequest) {
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
}

func (machineService *ServiceImpl) Delete(ginContext *gin.Context, deleteMachineRequest *model.DeleteMachineRequest) {
	valErr := machineService.validatorService.ValidateStruct(deleteMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *deleteMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := machineService.machineRepository.Delete(gormTransaction, deleteMachineRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
