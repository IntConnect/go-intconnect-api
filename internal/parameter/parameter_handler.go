package parameter

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
	parameterService Service
	viperConfig      *viper.Viper
}

func NewHandler(parameterService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		parameterService: parameterService,
		viperConfig:      viperConfig,
	}
}

func (parameterHandler *Handler) FindAllParameter(ginContext *gin.Context) {
	parameterResponses := parameterHandler.parameterService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Parameters has been fetched", parameterResponses))
}

func (parameterHandler *Handler) FindAllParameterPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	parameterResponses := parameterHandler.parameterService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Parameter has been fetched", parameterResponses))
}

func (parameterHandler *Handler) CreateParameter(ginContext *gin.Context) {
	var createParameterModel model.CreateParameterRequest
	err := ginContext.ShouldBind(&createParameterModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	modelFile, err := ginContext.FormFile("model")
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))

	parameterHandler.parameterService.Create(ginContext, &createParameterModel, modelFile)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Parameter has been created", nil))
}

func (parameterHandler *Handler) UpdateParameter(ginContext *gin.Context) {
	var updateParameterModel model.UpdateParameterRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateParameterModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	parameterHandler.parameterService.Update(ginContext, &updateParameterModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Parameter has been created", nil))
}

func (parameterHandler *Handler) DeleteParameter(ginContext *gin.Context) {
	var deleteBomModel model.DeleteParameterRequest
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deleteBomModel.Id = parsedBomId
	parameterHandler.parameterService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
