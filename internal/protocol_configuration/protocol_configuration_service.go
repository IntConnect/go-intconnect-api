package protocol_configuration

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createProtocolConfigurationRequest *model.CreateProtocolConfigurationRequest) []*model.ProtocolConfigurationResponse
	FindAll() []*model.ProtocolConfigurationResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.ProtocolConfigurationResponse]
	Update(ginContext *gin.Context, updateProtocolConfigurationRequest *model.UpdateProtocolConfigurationRequest)
	Delete(ginContext *gin.Context, deleteProtocolConfigurationRequest *model.DeleteProtocolConfigurationRequest)
	FindById(ginContext *gin.Context, protocolConfigurationId uint64) *model.ProtocolConfigurationResponse
}
