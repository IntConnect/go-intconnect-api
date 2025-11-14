package facility

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
	facilityRepository Repository
	validatorService   validator.Service
	dbConnection       *gorm.DB
	viperConfig        *viper.Viper
}

func NewService(facilityRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		facilityRepository: facilityRepository,
		validatorService:   validatorService,
		dbConnection:       dbConnection,
		viperConfig:        viperConfig,
	}
}

func (facilityService *ServiceImpl) FindAll() []*model.FacilityResponse {
	var facilityResponsesRequest []*model.FacilityResponse
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facilityEntities, err := facilityService.facilityRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		facilityResponsesRequest = mapper.MapFacilityEntitiesIntoFacilityResponses(facilityEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return facilityResponsesRequest
}

func (facilityService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.FacilityResponse] {
	paginationResp := model.PaginationResponse[*model.FacilityResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allFacility []*model.FacilityResponse
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facilityEntities, totalItems, err := facilityService.facilityRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allFacility = mapper.MapFacilityEntitiesIntoFacilityResponses(facilityEntities)
		paginationResp = model.PaginationResponse[*model.FacilityResponse]{
			Data:        allFacility,
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

// Create - Membuat facility baru
func (facilityService *ServiceImpl) Create(ginContext *gin.Context, createFacilityRequest *model.CreateFacilityRequest) {
	valErr := facilityService.validatorService.ValidateStruct(createFacilityRequest)
	facilityService.validatorService.ParseValidationError(valErr, *createFacilityRequest)
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facilityEntity := mapper.MapCreateFacilityRequestIntoFacilityEntity(createFacilityRequest)
		err := facilityService.facilityRepository.Create(gormTransaction, facilityEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (facilityService *ServiceImpl) Update(ginContext *gin.Context, updateFacilityRequest *model.UpdateFacilityRequest) {
	valErr := facilityService.validatorService.ValidateStruct(updateFacilityRequest)
	facilityService.validatorService.ParseValidationError(valErr, *updateFacilityRequest)
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facility, err := facilityService.facilityRepository.FindById(gormTransaction, updateFacilityRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdateFacilityRequestIntoFacilityEntity(updateFacilityRequest, facility)
		err = facilityService.facilityRepository.Update(gormTransaction, facility)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (facilityService *ServiceImpl) Delete(ginContext *gin.Context, deleteFacilityRequest *model.DeleteFacilityRequest) {
	valErr := facilityService.validatorService.ValidateStruct(deleteFacilityRequest)
	facilityService.validatorService.ParseValidationError(valErr, *deleteFacilityRequest)
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := facilityService.facilityRepository.Delete(gormTransaction, deleteFacilityRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
