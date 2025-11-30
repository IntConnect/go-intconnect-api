package node

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
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	nodeResponses := nodeHandler.nodeService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been fetched", nodeResponses))
}

func (nodeHandler *Handler) CreateNode(ginContext *gin.Context) {
	var createNodeModel model.CreateNodeRequest
	err := ginContext.ShouldBindBodyWithJSON(&createNodeModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	nodeHandler.nodeService.Create(ginContext, &createNodeModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been created", nil))
}

func (nodeHandler *Handler) UpdateNode(ginContext *gin.Context) {
	var updateNodeModel model.UpdateNodeRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateNodeModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	nodeHandler.nodeService.Update(ginContext, &updateNodeModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been created", nil))
}

func (nodeHandler *Handler) DeleteNode(ginContext *gin.Context) {
	var deleteNodeModel model.DeleteNodeRequest
	nodeId := ginContext.Param("id")
	err := ginContext.ShouldBindBodyWithJSON(&deleteNodeModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	parsedNodeId, err := strconv.ParseUint(nodeId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteNodeModel.Id = parsedNodeId
	nodeHandler.nodeService.Delete(ginContext, &deleteNodeModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Node has been updated", nil))
}
