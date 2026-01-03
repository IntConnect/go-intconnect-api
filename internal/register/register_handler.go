package register

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
	registerService Service
	viperConfig     *viper.Viper
}

func NewHandler(registerService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		registerService: registerService,
		viperConfig:     viperConfig,
	}
}

func (registerHandler *Handler) FindAllRegister(ginContext *gin.Context) {
	registerResponses := registerHandler.registerService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Register has been fetched", registerResponses))
}

func (registerHandler *Handler) FindAllRegisterPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := registerHandler.registerService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (registerHandler *Handler) FindRegisterDependency(ginContext *gin.Context) {
	mqttTopicDependency := registerHandler.registerService.FindDependency()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse(model.RESPONSE_SUCCESS, mqttTopicDependency))
}

func (registerHandler *Handler) CreateRegister(ginContext *gin.Context) {
	var createRegisterModel model.CreateRegisterRequest
	err := ginContext.ShouldBindBodyWithJSON(&createRegisterModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := registerHandler.registerService.Create(ginContext, &createRegisterModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (registerHandler *Handler) UpdateRegister(ginContext *gin.Context) {
	var updateRegisterModel model.UpdateRegisterRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateRegisterModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	registerId := ginContext.Param("id")
	parsedRegisterId, err := strconv.ParseUint(registerId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateRegisterModel.Id = parsedRegisterId
	paginatedResponse := registerHandler.registerService.Update(ginContext, &updateRegisterModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (registerHandler *Handler) DeleteRegister(ginContext *gin.Context) {
	var deleteRegisterModel model.DeleteResourceGeneralRequest
	registerId := ginContext.Param("id")
	parsedRegisterId, err := strconv.ParseUint(registerId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))

	err = ginContext.ShouldBindBodyWithJSON(&deleteRegisterModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteRegisterModel.Id = parsedRegisterId
	paginatedResponse := registerHandler.registerService.Delete(ginContext, &deleteRegisterModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
