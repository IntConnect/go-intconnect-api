package node

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"
)

type Handler struct {
	nodeService Service
	viperConfig *viper.Viper
}

func NewHandler(nodeService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		nodeService: nodeService,
		viperConfig: viperConfig,
	}
}

func (nodeHandler *Handler) FindAll(ginContext *gin.Context) {
	nodeResponses := nodeHandler.nodeService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been fetched", nodeResponses))
}

func (nodeHandler *Handler) FindAllPagination(ginContext *gin.Context) {
	paginationReq := model.PaginationRequest{
		Page:  1,
		Size:  10,
		Sort:  "id",
		Order: "asc",
	}

	// Bind query parameters to the struct
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	nodeResponses := nodeHandler.nodeService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been fetched", nodeResponses))
}

func (nodeHandler *Handler) CreateNode(ginContext *gin.Context) {
	var createNodeModel model.CreateNodeDto
	err := ginContext.ShouldBindBodyWithJSON(&createNodeModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	nodeHandler.nodeService.Create(ginContext, &createNodeModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been created", nil))
}

func (nodeHandler *Handler) UpdateNode(ginContext *gin.Context) {
	var updateNodeModel model.UpdateNodeDto
	err := ginContext.ShouldBindBodyWithJSON(&updateNodeModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	nodeHandler.nodeService.Update(ginContext, &updateNodeModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been created", nil))
}

func (nodeHandler *Handler) DeleteNode(ginContext *gin.Context) {
	var deleteBomModel model.DeleteNodeDto
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deleteBomModel.Id = parsedBomId
	nodeHandler.nodeService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
