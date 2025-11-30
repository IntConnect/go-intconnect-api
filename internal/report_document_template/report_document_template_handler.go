package report_document_template

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
	reportDocumentTemplateService Service
	viperConfig                   *viper.Viper
}

func NewHandler(reportDocumentTemplateService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		reportDocumentTemplateService: reportDocumentTemplateService,
		viperConfig:                   viperConfig,
	}
}

func (reportDocumentTemplateHandler *Handler) FindAllReportDocumentTemplate(ginContext *gin.Context) {
	reportDocumentTemplateResponses := reportDocumentTemplateHandler.reportDocumentTemplateService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Report document template has been fetched", reportDocumentTemplateResponses))
}

func (reportDocumentTemplateHandler *Handler) FindAllReportDocumentTemplatePagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := reportDocumentTemplateHandler.reportDocumentTemplateService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (reportDocumentTemplateHandler *Handler) CreateReportDocumentTemplate(ginContext *gin.Context) {
	var createReportDocumentTemplateModel model.CreateReportDocumentTemplateRequest
	err := ginContext.ShouldBindBodyWithJSON(&createReportDocumentTemplateModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	reportDocumentTemplateHandler.reportDocumentTemplateService.Create(ginContext, &createReportDocumentTemplateModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("ReportDocumentTemplate has been created", nil))
}

func (reportDocumentTemplateHandler *Handler) UpdateReportDocumentTemplate(ginContext *gin.Context) {
	var updateReportDocumentTemplateModel model.UpdateReportDocumentTemplateRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateReportDocumentTemplateModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	reportDocumentTemplateHandler.reportDocumentTemplateService.Update(ginContext, &updateReportDocumentTemplateModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("ReportDocumentTemplate has been created", nil))
}

func (reportDocumentTemplateHandler *Handler) DeleteReportDocumentTemplate(ginContext *gin.Context) {
	var deleteResourceGeneralRequest model.DeleteResourceGeneralRequest
	reportDocumentTemplateId := ginContext.Param("id")
	parsedReportDocumentTemplateId, err := strconv.ParseUint(reportDocumentTemplateId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteResourceGeneralRequest.Id = parsedReportDocumentTemplateId
	reportDocumentTemplateHandler.reportDocumentTemplateService.Delete(ginContext, &deleteResourceGeneralRequest)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("reportRocumentTemplateId has been updated", nil))
}
