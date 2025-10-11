package protocol_configuration

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createProtocolConfigurationDto *model.CreateProtocolConfigurationDto)
	FindAll() []*model.ProtocolConfigurationResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.ProtocolConfigurationResponse]
	Update(ginContext *gin.Context, updateProtocolConfigurationDto *model.UpdateProtocolConfigurationDto)
	Delete(ginContext *gin.Context, deleteProtocolConfigurationDto *model.DeleteProtocolConfigurationDto)
}
