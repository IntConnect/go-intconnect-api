package modbus_server

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.ModbusServerResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.ModbusServerResponse]
	Create(ginContext *gin.Context, createModbusServerRequest *model.CreateModbusServerRequest) *model.PaginatedResponse[*model.ModbusServerResponse]
	Update(ginContext *gin.Context, updateModbusServerRequest *model.UpdateModbusServerRequest) *model.PaginatedResponse[*model.ModbusServerResponse]
	Delete(ginContext *gin.Context, deleteModbusServerRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.ModbusServerResponse]
}
