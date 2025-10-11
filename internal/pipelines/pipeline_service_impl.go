package pipeline

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	pipelineRepository Repository
	validatorService   validator.Service
	dbConnection       *gorm.DB
	viperConfig        *viper.Viper
}

func NewService(pipelineRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		pipelineRepository: pipelineRepository,
		validatorService:   validatorService,
		dbConnection:       dbConnection,
		viperConfig:        viperConfig,
	}
}

func (pipelineService *ServiceImpl) FindAll() []*model.PipelineResponse {
	var pipelineResponsesDto []*model.PipelineResponse
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntities, err := pipelineService.pipelineRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		pipelineResponsesDto = mapper.MapPipelineEntitiesIntoPipelineResponses(pipelineEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return pipelineResponsesDto
}

func (pipelineService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.PipelineResponse] {
	paginationResp := model.PaginationResponse[*model.PipelineResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allPipeline []*model.PipelineResponse
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntities, totalItems, err := pipelineService.pipelineRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allPipeline = mapper.MapPipelineEntitiesIntoPipelineResponses(pipelineEntities)
		paginationResp = model.PaginationResponse[*model.PipelineResponse]{
			Data:        allPipeline,
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

// Create - Membuat pipeline baru
func (pipelineService *ServiceImpl) Create(ginContext *gin.Context, createPipelineDto *model.CreatePipelineDto) {
	valErr := pipelineService.validatorService.ValidateStruct(createPipelineDto)
	pipelineService.validatorService.ParseValidationError(valErr, *createPipelineDto)
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntity := mapper.MapCreatePipelineDtoIntoPipelineEntity(createPipelineDto)
		err := pipelineService.pipelineRepository.Create(gormTransaction, pipelineEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (pipelineService *ServiceImpl) Update(ginContext *gin.Context, updatePipelineDto *model.UpdatePipelineDto) {
	valErr := pipelineService.validatorService.ValidateStruct(updatePipelineDto)
	pipelineService.validatorService.ParseValidationError(valErr, *updatePipelineDto)
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipeline, err := pipelineService.pipelineRepository.FindById(gormTransaction, updatePipelineDto.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdatePipelineDtoIntoPipelineEntity(updatePipelineDto, pipeline)
		err = pipelineService.pipelineRepository.Update(gormTransaction, pipeline)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (pipelineService *ServiceImpl) Delete(ginContext *gin.Context, deletePipelineDto *model.DeletePipelineDto) {
	valErr := pipelineService.validatorService.ValidateStruct(deletePipelineDto)
	pipelineService.validatorService.ParseValidationError(valErr, *deletePipelineDto)
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := pipelineService.pipelineRepository.Delete(gormTransaction, deletePipelineDto.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
