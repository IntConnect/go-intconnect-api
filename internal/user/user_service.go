package user

import (
	"github.com/gin-gonic/gin"
	"go-intconnect-api/internal/model"
)

type Service interface {
	Create(ginContext *gin.Context, createUserDto *model.CreateUserDto)
	FindAll() []*model.UserResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.UserResponse]
	HandleLogin(ginContext *gin.Context, loginUserDto *model.LoginUserDto) string
	Update(ginContext *gin.Context, updateUserDto *model.UpdateUserDto)
	Delete(ginContext *gin.Context, deleteUserDto *model.DeleteUserDto)
}
