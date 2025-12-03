package machine

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/spf13/viper"
)

type Handler struct {
	machineService Service
	formDecoder    *form.Decoder
	viperConfig    *viper.Viper
}

func NewHandler(machineService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		machineService: machineService,
		viperConfig:    viperConfig,
		formDecoder:    form.NewDecoder(),
	}
}

func (machineHandler *Handler) FindAllMachine(ginContext *gin.Context) {
	machineResponses := machineHandler.machineService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Machine has been fetched", machineResponses))
}

func (machineHandler *Handler) FindAllMachinePagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := machineHandler.machineService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (machineHandler *Handler) CreateMachine(ginContext *gin.Context) {

	var createMachineModel model.CreateMachineRequest
	err := ginContext.Request.ParseMultipartForm(32 << 20) // 32MB maxMemory
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	err = machineHandler.formDecoder.Decode(&createMachineModel, ginContext.Request.PostForm)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	modelFile, _ := ginContext.FormFile("model")
	thumbnailFile, _ := ginContext.FormFile("thumbnail")
	createMachineModel.Model = modelFile
	createMachineModel.Thumbnail = thumbnailFile
	extractIndexedFiles, err := helper.ExtractIndexedFiles(ginContext, "machine_documents[", "].document_file", len(createMachineModel.MachineDocuments))
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	for i, machineDocument := range createMachineModel.MachineDocuments {
		machineDocument.DocumentFile = extractIndexedFiles[i]
		createMachineModel.MachineDocuments[i] = machineDocument
	}
	if len(createMachineModel.MachineDocuments) > 0 {
		machineHandler.machineService.Create(ginContext, &createMachineModel)
	}
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Machine has been created", nil))
}

func (machineHandler *Handler) UpdateMachine(ginContext *gin.Context) {
	var updateMachineModel model.UpdateMachineRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateMachineModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	machineHandler.machineService.Update(ginContext, &updateMachineModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Machine has been created", nil))
}

func (machineHandler *Handler) DeleteMachine(ginContext *gin.Context) {
	var deleteBomModel model.DeleteMachineRequest
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteBomModel.Id = parsedBomId
	machineHandler.machineService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
