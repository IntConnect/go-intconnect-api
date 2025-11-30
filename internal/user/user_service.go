package user

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.UserResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.UserResponse]
	HandleLogin(ginContext *gin.Context, loginUserRequest *model.LoginUserRequest) string
	Create(ginContext *gin.Context, createUserRequest *model.CreateUserRequest) *model.PaginatedResponse[*model.UserResponse]
	Update(ginContext *gin.Context, updateUserRequest *model.UpdateUserRequest) *model.PaginatedResponse[*model.UserResponse]
	Delete(ginContext *gin.Context, deleteUserRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.UserResponse]
}
