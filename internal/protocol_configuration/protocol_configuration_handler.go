package protocol_configuration

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	protocolConfigurationService Service
	viperConfig                  *viper.Viper
}

func NewHandler(protocolConfigurationService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		protocolConfigurationService: protocolConfigurationService,
		viperConfig:                  viperConfig,
	}
}

func (protocolConfigurationHandler *Handler) FindAll(ginContext *gin.Context) {
	protocolConfigurationResponses := protocolConfigurationHandler.protocolConfigurationService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("ProtocolConfiguration has been fetched", protocolConfigurationResponses))
}

func (protocolConfigurationHandler *Handler) FindById(ginContext *gin.Context) {
	protocolConfigurationId := ginContext.Param("id")
	parsedProtocolConfigurationId, err := strconv.ParseUint(protocolConfigurationId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, "Supplied param not valid", nil))
	protocolConfigurationResponses := protocolConfigurationHandler.protocolConfigurationService.FindById(ginContext, parsedProtocolConfigurationId)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("ProtocolConfiguration has been fetched", protocolConfigurationResponses))
}

func (protocolConfigurationHandler *Handler) FindAllPagination(ginContext *gin.Context) {
	paginationReq := model.PaginationRequest{
		Page:  1,
		Size:  10,
		Sort:  "id",
		Order: "asc",
	}

	// Bind query parameters to the struct
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	protocolConfigurationResponses := protocolConfigurationHandler.protocolConfigurationService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("ProtocolConfiguration has been fetched", protocolConfigurationResponses))
}

func (protocolConfigurationHandler *Handler) CreateProtocolConfiguration(ginContext *gin.Context) {
	var createProtocolConfigurationModel model.CreateProtocolConfigurationRequest
	err := ginContext.ShouldBindBodyWithJSON(&createProtocolConfigurationModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	protocolConfigurationResponses := protocolConfigurationHandler.protocolConfigurationService.Create(ginContext, &createProtocolConfigurationModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("ProtocolConfiguration has been created", protocolConfigurationResponses))
}

func (protocolConfigurationHandler *Handler) UpdateProtocolConfiguration(ginContext *gin.Context) {
	var updateProtocolConfigurationModel model.UpdateProtocolConfigurationRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateProtocolConfigurationModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	protocolConfigurationHandler.protocolConfigurationService.Update(ginContext, &updateProtocolConfigurationModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("ProtocolConfiguration has been created", nil))
}

func (protocolConfigurationHandler *Handler) DeleteProtocolConfiguration(ginContext *gin.Context) {
	var deleteBomModel model.DeleteProtocolConfigurationRequest
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deleteBomModel.Id = parsedBomId
	protocolConfigurationHandler.protocolConfigurationService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
