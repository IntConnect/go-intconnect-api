package user

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
	userService Service
	viperConfig *viper.Viper
}

func NewHandler(userService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		userService: userService,
		viperConfig: viperConfig,
	}
}

func (userHandler *Handler) FindAllUser(ginContext *gin.Context) {
	userResponses := userHandler.userService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("User has been fetched", userResponses))
}

func (userHandler *Handler) FindAllUserPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	userResponses := userHandler.userService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("User has been fetched", userResponses))
}

func (userHandler *Handler) LoginUser(ginContext *gin.Context) {
	var loginUserRequest model.LoginUserRequest
	err := ginContext.ShouldBindBodyWithJSON(&loginUserRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	generatedToken := userHandler.userService.HandleLogin(ginContext, &loginUserRequest)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("User logged successfully", gin.H{
		"token": generatedToken,
	}))
}

func (userHandler *Handler) CreateUser(ginContext *gin.Context) {
	var createUserModel model.CreateUserRequest
	err := ginContext.ShouldBindBodyWithJSON(&createUserModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	userHandler.userService.Create(ginContext, &createUserModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("User has been created", nil))
}

func (userHandler *Handler) UpdateUser(ginContext *gin.Context) {
	var updateUserModel model.UpdateUserRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateUserModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	userHandler.userService.Update(ginContext, &updateUserModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("User has been created", nil))
}

func (userHandler *Handler) DeleteUser(ginContext *gin.Context) {
	var deleteBomModel model.DeleteUserRequest
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deleteBomModel.Id = parsedBomId
	userHandler.userService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
