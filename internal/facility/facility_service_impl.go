package facility

import (
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/storage"
	"go-intconnect-api/internal/trait"
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
	facilityRepository  Repository
	auditLogService     auditLog.Service
	validatorService    validator.Service
	dbConnection        *gorm.DB
	viperConfig         *viper.Viper
	localStorageService *storage.Manager
}

func NewService(facilityRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	localStorageService *storage.Manager,
	auditLogService auditLog.Service) *ServiceImpl {
	return &ServiceImpl{
		facilityRepository:  facilityRepository,
		validatorService:    validatorService,
		dbConnection:        dbConnection,
		viperConfig:         viperConfig,
		localStorageService: localStorageService,
		auditLogService:     auditLogService,
	}
}

func (facilityService *ServiceImpl) FindAll() []*model.FacilityResponse {
	var facilityResponsesRequest []*model.FacilityResponse
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facilityEntities, err := facilityService.facilityRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		facilityResponsesRequest = helper.MapEntitiesIntoResponsesWithFunc[*entity.Facility, *model.FacilityResponse](facilityEntities, mapper.FuncMapAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return facilityResponsesRequest
}

func (facilityService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.FacilityResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var facilityResponses []*model.FacilityResponse
	var totalItems int64

	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facilityEntities, total, err := facilityService.facilityRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		facilityResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.Facility, *model.FacilityResponse](
			facilityEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Facilities fetched successfully",
		facilityResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat facility baru
func (facilityService *ServiceImpl) Create(ginContext *gin.Context, createFacilityRequest *model.CreateFacilityRequest) *model.PaginatedResponse[*model.FacilityResponse] {
	var paginationResp *model.PaginatedResponse[*model.FacilityResponse]
	valErr := facilityService.validatorService.ValidateStruct(createFacilityRequest)
	facilityService.validatorService.ParseValidationError(valErr, *createFacilityRequest)
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facilityEntity := helper.MapCreateRequestIntoEntity[model.CreateFacilityRequest, entity.Facility](createFacilityRequest)
		facilityEntity.Status = trait.FacilityStatusActive
		thumbnailPath, err := facilityService.localStorageService.Disk().Put(createFacilityRequest.Thumbnail, fmt.Sprintf("facilities/thumbnails/%d-%s", time.Now().UnixNano(), createFacilityRequest.Thumbnail.Filename))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
		facilityEntity.ThumbnailPath = thumbnailPath
		err = facilityService.facilityRepository.Create(gormTransaction, facilityEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = facilityService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (facilityService *ServiceImpl) Update(
	ginContext *gin.Context,
	updateFacilityRequest *model.UpdateFacilityRequest,
) *model.PaginatedResponse[*model.FacilityResponse] {
	userJwtClaim := helper.ExtractJwtClaimFromContext(ginContext)
	ipAddress, _ := helper.ExtractRequestMeta(ginContext)
	var paginationResp *model.PaginatedResponse[*model.FacilityResponse]

	valErr := facilityService.validatorService.ValidateStruct(updateFacilityRequest)
	facilityService.validatorService.ParseValidationError(valErr, *updateFacilityRequest)

	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {

		// 1. Ambil entity lama
		facility, err := facilityService.facilityRepository.FindById(gormTransaction, updateFacilityRequest.Id)
		pastFacility := facility
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		// 2. Update fields non-Thumbnail
		helper.MapUpdateRequestIntoEntity(updateFacilityRequest, facility)

		// 3. Cek apakah ada upload Thumbnail baru
		if updateFacilityRequest.Thumbnail != nil {
			// 3a. Hapus file lama jika ada
			if facility.ThumbnailPath != "" {
				_ = facilityService.localStorageService.Disk().Delete(facility.ThumbnailPath)
			}

			// 3b. Simpan file baru
			newPath, err := facilityService.localStorageService.Disk().Put(
				updateFacilityRequest.Thumbnail,
				fmt.Sprintf("facilities/thumbnails/%d-%s", time.Now().UnixNano(), updateFacilityRequest.Thumbnail.Filename),
			)
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))

			facility.ThumbnailPath = newPath
		}

		facility.Auditable = entity.UpdateAuditable(userJwtClaim.Name)
		err = facilityService.facilityRepository.Update(gormTransaction, facility)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		facilityService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      userJwtClaim.Id,
			Action:      model.AUDIT_LOG_CREATE,
			Feature:     model.AUDIT_LOG_FEATURE_FACILITY,
			Description: "",
			Before:      pastFacility,
			After:       facility,
			IpAddress:   ipAddress,
		})
		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))

	// Kembalikan pagination
	paginationRequest := model.NewPaginationRequest()
	paginationResp = facilityService.FindAllPagination(&paginationRequest)

	return paginationResp
}

func (facilityService *ServiceImpl) Delete(ginContext *gin.Context, deleteFacilityRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.FacilityResponse] {
	userJwtClaim := helper.ExtractJwtClaimFromContext(ginContext)
	ipAddress, _ := helper.ExtractRequestMeta(ginContext)

	var paginationResp *model.PaginatedResponse[*model.FacilityResponse]
	valErr := facilityService.validatorService.ValidateStruct(deleteFacilityRequest)
	facilityService.validatorService.ParseValidationError(valErr, *deleteFacilityRequest)
	err := facilityService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		facilityEntity, err := facilityService.facilityRepository.FindById(gormTransaction, deleteFacilityRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		if facilityEntity.ThumbnailPath != "" {
			_ = facilityService.localStorageService.Disk().Delete(facilityEntity.ThumbnailPath)
		}
		err = facilityService.facilityRepository.Delete(gormTransaction, deleteFacilityRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		facilityService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      userJwtClaim.Id,
			Action:      model.AUDIT_LOG_CREATE,
			Feature:     model.AUDIT_LOG_FEATURE_FACILITY,
			Description: deleteFacilityRequest.Reason,
			Before:      facilityEntity,
			After:       nil,
			IpAddress:   ipAddress,
		})
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = facilityService.FindAllPagination(&paginationRequest)
	return paginationResp
}
