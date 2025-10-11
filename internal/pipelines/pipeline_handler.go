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
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been fetched", pipelineResponses))
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
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been fetched", pipelineResponses))
}

func (pipelineHandler *Handler) CreateNode(ginContext *gin.Context) {
	var createNodeModel model.CreateNodeDto
	err := ginContext.ShouldBindBodyWithJSON(&createNodeModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineHandler.pipelineService.Create(ginContext, &createNodeModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been created", nil))
}

func (pipelineHandler *Handler) UpdateNode(ginContext *gin.Context) {
	var updateNodeModel model.UpdateNodeDto
	err := ginContext.ShouldBindBodyWithJSON(&updateNodeModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineHandler.pipelineService.Update(ginContext, &updateNodeModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been created", nil))
}

func (pipelineHandler *Handler) DeleteNode(ginContext *gin.Context) {
	var deleteBomModel model.DeleteNodeDto
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deleteBomModel.Id = parsedBomId
	pipelineHandler.pipelineService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
