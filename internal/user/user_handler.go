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

func (userHandler *Handler) FindById(ginContext *gin.Context) {
	userId := ginContext.Param("id")
	parsedUserId, err := strconv.ParseUint(userId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrParameterInvalid))
	userResponse := userHandler.userService.FindById(ginContext, parsedUserId)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("User fetched successfully", userResponse))
}

func (userHandler *Handler) FindAllUserPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := userHandler.userService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (userHandler *Handler) LoginUser(ginContext *gin.Context) {
	var loginUserRequest model.LoginUserRequest
	err := ginContext.ShouldBindBodyWithJSON(&loginUserRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	generatedToken := userHandler.userService.HandleLogin(ginContext, &loginUserRequest)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("User logged successfully", gin.H{
		"token": generatedToken,
	}))
}

func (userHandler *Handler) CreateUser(ginContext *gin.Context) {
	var createUserModel model.CreateUserRequest
	err := ginContext.ShouldBindBodyWithJSON(&createUserModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := userHandler.userService.Create(ginContext, &createUserModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (userHandler *Handler) UpdateUser(ginContext *gin.Context) {
	var updateUserModel model.UpdateUserRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateUserModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	userId := ginContext.Param("id")
	parsedUserId, err := strconv.ParseUint(userId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateUserModel.Id = parsedUserId
	paginatedResponse := userHandler.userService.Update(ginContext, &updateUserModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (userHandler *Handler) DeleteUser(ginContext *gin.Context) {
	var deleteUserModel model.DeleteResourceGeneralRequest
	err := ginContext.ShouldBindBodyWithJSON(&deleteUserModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	userId := ginContext.Param("id")
	parsedUserId, err := strconv.ParseUint(userId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteUserModel.Id = parsedUserId
	paginatedResponse := userHandler.userService.Delete(ginContext, &deleteUserModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
