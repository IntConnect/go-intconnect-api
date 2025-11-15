package machine

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"math"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	machineRepository Repository
	validatorService  validator.Service
	dbConnection      *gorm.DB
	viperConfig       *viper.Viper
}

func NewService(machineRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		machineRepository: machineRepository,
		validatorService:  validatorService,
		dbConnection:      dbConnection,
		viperConfig:       viperConfig,
	}
}

func (machineService *ServiceImpl) FindAll() []*model.MachineResponse {
	var machineResponsesRequest []*model.MachineResponse
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntities, err := machineService.machineRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		machineResponsesRequest = helper.MapEntitiesIntoResponses[entity.Machine, model.MachineResponse](machineEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return machineResponsesRequest
}

func (machineService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.MachineResponse] {
	paginationResp := model.PaginationResponse[*model.MachineResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allMachine []*model.MachineResponse
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntities, totalItems, err := machineService.machineRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allMachine = helper.MapEntitiesIntoResponses[entity.Machine, model.MachineResponse](machineEntities)
		paginationResp = model.PaginationResponse[*model.MachineResponse]{
			Data:        allMachine,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: paginationReq.Page,
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}

// Create - Membuat machine baru
func (machineService *ServiceImpl) Create(ginContext *gin.Context, createMachineRequest *model.CreateMachineRequest, modelFile *multipart.FileHeader) {
	valErr := machineService.validatorService.ValidateStruct(createMachineRequest)
	machineService.validatorService.ParseValidationError(valErr, *createMachineRequest)
	err := machineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		machineEntity := helper.MapCreateRequestIntoEntity[model.CreateMachineRequest, entity.Machine](createMachineRequest)
		err := machineService.machineRepository.Create(gormTransaction, machineEntity)
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
