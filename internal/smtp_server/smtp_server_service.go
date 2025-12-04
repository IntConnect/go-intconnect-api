package smtp_server

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.SmtpServerResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.SmtpServerResponse]
	Create(ginContext *gin.Context, createSmtpServerRequest *model.CreateSmtpServerRequest) *model.PaginatedResponse[*model.SmtpServerResponse]
	Update(ginContext *gin.Context, updateSmtpServerRequest *model.UpdateSmtpServerRequest) *model.PaginatedResponse[*model.SmtpServerResponse]
	Delete(ginContext *gin.Context, deleteSmtpServerRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.SmtpServerResponse]
}
