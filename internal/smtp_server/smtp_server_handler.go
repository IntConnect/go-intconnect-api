package smtp_server

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
	smtpServerService Service
	viperConfig       *viper.Viper
}

func NewHandler(smtpServerService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		smtpServerService: smtpServerService,
		viperConfig:       viperConfig,
	}
}

func (smtpServerHandler *Handler) FindAllSmtpServer(ginContext *gin.Context) {
	smtpServerResponses := smtpServerHandler.smtpServerService.FindAll()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("SMTP Server has been fetched", smtpServerResponses))
}

func (smtpServerHandler *Handler) FindAllSmtpServerPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := smtpServerHandler.smtpServerService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (smtpServerHandler *Handler) CreateSmtpServer(ginContext *gin.Context) {
	var createSmtpServerModel model.CreateSmtpServerRequest
	err := ginContext.ShouldBindBodyWithJSON(&createSmtpServerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := smtpServerHandler.smtpServerService.Create(ginContext, &createSmtpServerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (smtpServerHandler *Handler) UpdateSmtpServer(ginContext *gin.Context) {
	var updateSmtpServerModel model.UpdateSmtpServerRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateSmtpServerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	smtpServerId := ginContext.Param("id")
	parsedSmtpServerId, err := strconv.ParseUint(smtpServerId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateSmtpServerModel.Id = parsedSmtpServerId
	paginatedResponse := smtpServerHandler.smtpServerService.Update(ginContext, &updateSmtpServerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (smtpServerHandler *Handler) DeleteSmtpServer(ginContext *gin.Context) {
	var deleteSmtpServerModel model.DeleteResourceGeneralRequest
	smtpServerId := ginContext.Param("id")
	parsedSmtpServerId, err := strconv.ParseUint(smtpServerId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))

	err = ginContext.ShouldBindBodyWithJSON(&deleteSmtpServerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteSmtpServerModel.Id = parsedSmtpServerId
	paginatedResponse := smtpServerHandler.smtpServerService.Delete(ginContext, &deleteSmtpServerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
