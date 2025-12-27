package check_sheet_document_template

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
	checkSheetDocumentTemplateService Service
	viperConfig                       *viper.Viper
}

func NewHandler(checkSheetDocumentTemplateService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		checkSheetDocumentTemplateService: checkSheetDocumentTemplateService,
		viperConfig:                       viperConfig,
	}
}

func (checkSheetDocumentTemplateHandler *Handler) FindAllCheckSheetDocumentTemplate(ginContext *gin.Context) {
	checkSheetDocumentTemplateResponses := checkSheetDocumentTemplateHandler.checkSheetDocumentTemplateService.FindAll()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Check sheet document template has been fetched", checkSheetDocumentTemplateResponses))
}

func (checkSheetDocumentTemplateHandler *Handler) FindAllCheckSheetDocumentTemplatePagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := checkSheetDocumentTemplateHandler.checkSheetDocumentTemplateService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (checkSheetDocumentTemplateHandler *Handler) CreateCheckSheetDocumentTemplate(ginContext *gin.Context) {
	var createCheckSheetDocumentTemplateModel model.CreateCheckSheetDocumentTemplateRequest
	err := ginContext.ShouldBindBodyWithJSON(&createCheckSheetDocumentTemplateModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedRes := checkSheetDocumentTemplateHandler.checkSheetDocumentTemplateService.Create(ginContext, &createCheckSheetDocumentTemplateModel)
	ginContext.JSON(http.StatusOK, paginatedRes)
}

func (checkSheetDocumentTemplateHandler *Handler) UpdateCheckSheetDocumentTemplate(ginContext *gin.Context) {
	var updateCheckSheetDocumentTemplateModel model.UpdateCheckSheetDocumentTemplateRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateCheckSheetDocumentTemplateModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrPayloadInvalid))
	checkSheetDocumentTemplateId := ginContext.Param("id")
	parsedCheckSheetDocumentTemplateId, err := strconv.ParseUint(checkSheetDocumentTemplateId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrParameterInvalid))
	updateCheckSheetDocumentTemplateModel.Id = parsedCheckSheetDocumentTemplateId
	paginatedRes := checkSheetDocumentTemplateHandler.checkSheetDocumentTemplateService.Update(ginContext, &updateCheckSheetDocumentTemplateModel)
	ginContext.JSON(http.StatusOK, paginatedRes)
}

func (checkSheetDocumentTemplateHandler *Handler) DeleteCheckSheetDocumentTemplate(ginContext *gin.Context) {
	var deleteResourceGeneralRequest model.DeleteResourceGeneralRequest
	err := ginContext.ShouldBindBodyWithJSON(&deleteResourceGeneralRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrPayloadInvalid))
	checkSheetDocumentTemplateId := ginContext.Param("id")
	parsedCheckSheetDocumentTemplateId, err := strconv.ParseUint(checkSheetDocumentTemplateId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteResourceGeneralRequest.Id = parsedCheckSheetDocumentTemplateId
	paginatedRes := checkSheetDocumentTemplateHandler.checkSheetDocumentTemplateService.Delete(ginContext, &deleteResourceGeneralRequest)
	ginContext.JSON(http.StatusOK, paginatedRes)
}
