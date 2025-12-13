package breakdown

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
	breakdownService Service
	viperConfig      *viper.Viper
}

func NewHandler(breakdownService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		breakdownService: breakdownService,
		viperConfig:      viperConfig,
	}
}

func (breakdownHandler *Handler) FindAllBreakdown(ginContext *gin.Context) {
	breakdownResponses := breakdownHandler.breakdownService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("MQTT Topic has been fetched", breakdownResponses))
}

func (breakdownHandler *Handler) FindAllBreakdownPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := breakdownHandler.breakdownService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (breakdownHandler *Handler) FindBreakdownById(ginContext *gin.Context) {
	breakdownId := ginContext.Param("id")
	parsedBreakdownId, err := strconv.ParseUint(breakdownId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := breakdownHandler.breakdownService.FindById(ginContext, parsedBreakdownId)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (breakdownHandler *Handler) CreateBreakdown(ginContext *gin.Context) {
	var createBreakdownModel model.CreateBreakdownRequest
	err := ginContext.ShouldBindBodyWithJSON(&createBreakdownModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := breakdownHandler.breakdownService.Create(ginContext, &createBreakdownModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (breakdownHandler *Handler) UpdateBreakdown(ginContext *gin.Context) {
	var updateBreakdownModel model.UpdateBreakdownRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateBreakdownModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	breakdownId := ginContext.Param("id")
	parsedBreakdownId, err := strconv.ParseUint(breakdownId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateBreakdownModel.Id = parsedBreakdownId
	paginatedResponse := breakdownHandler.breakdownService.Update(ginContext, &updateBreakdownModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (breakdownHandler *Handler) DeleteBreakdown(ginContext *gin.Context) {
	var deleteBreakdownModel model.DeleteResourceGeneralRequest
	err := ginContext.ShouldBindBodyWithJSON(&deleteBreakdownModel)
	breakdownId := ginContext.Param("id")
	parsedBreakdownId, err := strconv.ParseUint(breakdownId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteBreakdownModel.Id = parsedBreakdownId
	paginatedResponse := breakdownHandler.breakdownService.Delete(ginContext, &deleteBreakdownModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
