package user

import (
	"github.com/gin-gonic/gin"
	"go-intconnect-api/internal/model"
)

type Service interface {
	Create(ginContext *gin.Context, createUserRequest *model.CreateUserRequest)
	FindAll() []*model.UserResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.UserResponse]
	HandleLogin(ginContext *gin.Context, loginUserRequest *model.LoginUserRequest) string
	Update(ginContext *gin.Context, updateUserRequest *model.UpdateUserRequest)
	Delete(ginContext *gin.Context, deleteUserRequest *model.DeleteUserRequest)
}
