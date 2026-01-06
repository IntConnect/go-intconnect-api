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

func (parameterHandler *Handler) FindDependencyParameter(ginContext *gin.Context) {
	parameterResponses := parameterHandler.parameterService.FindDependencyParameter()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("Parameters has been fetched", parameterResponses))
}

func (parameterHandler *Handler) FindAllParameter(ginContext *gin.Context) {
	var parameterFilterRequest model.ParameterFilterRequest
	err := ginContext.ShouldBindQuery(&parameterFilterRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	parameterResponses := parameterHandler.parameterService.FindAll(&parameterFilterRequest)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Parameters has been fetched", parameterResponses))
}

func (parameterHandler *Handler) FindAllParameterPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponses := parameterHandler.parameterService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponses)
}

func (parameterHandler *Handler) FindByIdParameter(ginContext *gin.Context) {
	parameterId := ginContext.Param("id")
	parsedParameterId, err := strconv.ParseUint(parameterId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	parameterResponses := parameterHandler.parameterService.FindById(ginContext, parsedParameterId)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("Parameters has been fetched", parameterResponses))
}

func (parameterHandler *Handler) CreateParameter(ginContext *gin.Context) {
	var createParameterRequest model.CreateParameterRequest
	err := ginContext.ShouldBindJSON(&createParameterRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	parameterHandler.parameterService.Create(ginContext, &createParameterRequest)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Parameter has been created", nil))
}

func (parameterHandler *Handler) UpdateParameter(ginContext *gin.Context) {
	var updateParameterRequest model.UpdateParameterRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateParameterRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	parameterId := ginContext.Param("id")
	parsedParameterId, err := strconv.ParseUint(parameterId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrParameterInvalid))
	updateParameterRequest.Id = parsedParameterId
	parameterHandler.parameterService.Update(ginContext, &updateParameterRequest)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Parameter has been updated", nil))
}

func (parameterHandler *Handler) UpdateParameterOperation(ginContext *gin.Context) {
	var manageParameterOperationModel model.ManageParameterOperationRequest
	err := ginContext.ShouldBindBodyWithJSON(&manageParameterOperationModel)
	parameterId := ginContext.Param("id")
	parsedParameterId, err := strconv.ParseUint(parameterId, 10, 64)
	manageParameterOperationModel.Id = parsedParameterId
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedRes := parameterHandler.parameterService.UpdateOperation(ginContext, &manageParameterOperationModel)
	ginContext.JSON(http.StatusOK, paginatedRes)
}

func (parameterHandler *Handler) DeleteParameter(ginContext *gin.Context) {
	var deleteParameterModel model.DeleteResourceGeneralRequest
	parameterId := ginContext.Param("id")
	parsedParameterId, err := strconv.ParseUint(parameterId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteParameterModel.Id = parsedParameterId
	parameterHandler.parameterService.Delete(ginContext, &deleteParameterModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Parameter has been deleted", nil))
}
