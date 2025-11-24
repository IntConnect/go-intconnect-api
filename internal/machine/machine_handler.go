package machine

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
	machineService Service
	viperConfig    *viper.Viper
}

func NewHandler(machineService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		machineService: machineService,
		viperConfig:    viperConfig,
	}
}

func (machineHandler *Handler) FindAllMachine(ginContext *gin.Context) {
	machineResponses := machineHandler.machineService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Machine has been fetched", machineResponses))
}

func (machineHandler *Handler) FindAllMachinePagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	paginatedResponse := machineHandler.machineService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (machineHandler *Handler) CreateMachine(ginContext *gin.Context) {
	var createMachineModel model.CreateMachineRequest

	err := ginContext.ShouldBind(&createMachineModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	modelFile, err := ginContext.FormFile("model")
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	createMachineModel.ModelHeader = modelFile
	machineHandler.machineService.Create(ginContext, &createMachineModel, modelFile)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Machine has been created", nil))
}

func (machineHandler *Handler) UpdateMachine(ginContext *gin.Context) {
	var updateMachineModel model.UpdateMachineRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateMachineModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	machineHandler.machineService.Update(ginContext, &updateMachineModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Machine has been created", nil))
}

func (machineHandler *Handler) DeleteMachine(ginContext *gin.Context) {
	var deleteBomModel model.DeleteMachineRequest
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deleteBomModel.Id = parsedBomId
	machineHandler.machineService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
