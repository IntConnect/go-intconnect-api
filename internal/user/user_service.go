package user

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.UserResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.UserResponse]
	HandleLogin(ginContext *gin.Context, loginUserRequest *model.LoginUserRequest) string
	Create(ginContext *gin.Context, createUserRequest *model.CreateUserRequest) model.PaginationResponse[*model.UserResponse]
	Update(ginContext *gin.Context, updateUserRequest *model.UpdateUserRequest)
	Delete(ginContext *gin.Context, deleteUserRequest *model.DeleteUserRequest)
}
