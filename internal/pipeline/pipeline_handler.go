package pipeline

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
	pipelineService Service
	viperConfig     *viper.Viper
}

func NewHandler(pipelineService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		pipelineService: pipelineService,
		viperConfig:     viperConfig,
	}
}

func (pipelineHandler *Handler) FindAll(ginContext *gin.Context) {
	pipelineResponses := pipelineHandler.pipelineService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Pipeline has been fetched", pipelineResponses))
}

func (pipelineHandler *Handler) FindById(ginContext *gin.Context) {
	pipelineId := ginContext.Param("id")
	parsedPipelineId, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
	}
	pipelineResponses := pipelineHandler.pipelineService.FindById(ginContext, parsedPipelineId)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Pipeline has been fetched", pipelineResponses))
}

func (pipelineHandler *Handler) RunPipeline(ginContext *gin.Context) {
	pipelineId := ginContext.Param("id")
	parsedPipelineId, err := strconv.ParseUint(pipelineId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, "Supplied param not valid", nil))
	pipelineResponses := pipelineHandler.pipelineService.RunPipeline(ginContext, parsedPipelineId)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Pipeline has been fetched", pipelineResponses))
}

func (pipelineHandler *Handler) FindAllPagination(ginContext *gin.Context) {
	paginationReq := model.PaginationRequest{
		Page:  1,
		Size:  10,
		Sort:  "id",
		Order: "asc",
	}

	// Bind query parameters to the struct
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineResponses := pipelineHandler.pipelineService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Pipeline has been fetched", pipelineResponses))
}

func (pipelineHandler *Handler) CreatePipeline(ginContext *gin.Context) {
	var createPipelineModel model.CreatePipelineRequest
	err := ginContext.ShouldBindBodyWithJSON(&createPipelineModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineHandler.pipelineService.Create(ginContext, &createPipelineModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Pipeline has been created", nil))
}

func (pipelineHandler *Handler) UpdatePipeline(ginContext *gin.Context) {
	var updatePipelineModel model.UpdatePipelineRequest
	err := ginContext.ShouldBindBodyWithJSON(&updatePipelineModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineHandler.pipelineService.Update(ginContext, &updatePipelineModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Pipeline has been created", nil))
}

func (pipelineHandler *Handler) DeletePipeline(ginContext *gin.Context) {
	var deletePipelineModel model.DeletePipelineRequest
	currencyId := ginContext.Param("id")
	parsedPipelineId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deletePipelineModel.ID = parsedPipelineId
	pipelineHandler.pipelineService.Delete(ginContext, &deletePipelineModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Pipeline has been updated", nil))
}
