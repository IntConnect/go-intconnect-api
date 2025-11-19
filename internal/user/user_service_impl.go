package user

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/trait"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	userRepository   Repository
	validatorService validator.Service
	dbConnection     *gorm.DB
	viperConfig      *viper.Viper
	loggerInstance   *logrus.Logger
}

func NewService(userRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, loggerInstance *logrus.Logger) *ServiceImpl {
	return &ServiceImpl{
		userRepository:   userRepository,
		validatorService: validatorService,
		dbConnection:     dbConnection,
		viperConfig:      viperConfig,
		loggerInstance:   loggerInstance,
	}
}

// Create - Membuat user baru
func (userService *ServiceImpl) FindAll() []*model.UserResponse {
	var allUser []*model.UserResponse
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userResponse, err := userService.userRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allUser = mapper.MapUserEntitiesIntoUserResponses(userResponse)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allUser
}

// Create - Membuat user baru
func (userService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.UserResponse] {
	paginationResp := model.PaginationResponse[*model.UserResponse]{}
	paginationQuery := helper.BuildPaginationQuery(paginationReq)

	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntities, totalItems, err := userService.userRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allUser := mapper.MapUserEntitiesIntoUserResponses(userEntities)
		paginationResp = model.PaginationResponse[*model.UserResponse]{
			Data:        allUser,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: paginationReq.Page,
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}

// Create - Membuat user baru
func (userService *ServiceImpl) Create(ginContext *gin.Context, createUserRequest *model.CreateUserRequest) model.PaginationResponse[*model.UserResponse] {
	paginationResp := model.PaginationResponse[*model.UserResponse]{}
	valErr := userService.validatorService.ValidateStruct(createUserRequest)
	userService.validatorService.ParseValidationError(valErr, *createUserRequest)
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntity := mapper.MapCreateUserRequestIntoUserEntity(createUserRequest)
		userEntity.Status = trait.UserStatusActive
		// TODO: Change it with logged user
		userEntity.Auditable = entity.NewAuditable("Administrator")
		err := userService.userRepository.Create(gormTransaction, userEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		paginationResp = userService.FindAllPagination(&paginationRequest)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}

func (userService *ServiceImpl) HandleLogin(ginContext *gin.Context, loginUserRequest *model.LoginUserRequest) string {
	err := userService.validatorService.ValidateStruct(loginUserRequest)
	userService.validatorService.ParseValidationError(err, loginUserRequest)
	var tokenString string
	err = userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userEntity entity.User
		err = gormTransaction.
			Where("email = ?", loginUserRequest.UserIdentifier).
			Or("username = ?", loginUserRequest.UserIdentifier).
			First(&userEntity).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		if err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(loginUserRequest.Password)); err != nil {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, "User credentials invalid", err))
		}
		tokenInstance := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email":    userEntity.Email,
			"username": userEntity.Username,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		})
		tokenString, err = tokenInstance.SignedString([]byte(userService.viperConfig.GetString("JWT_SECRET")))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError, err))
		if claims, ok := tokenInstance.Claims.(jwt.MapClaims); ok && tokenInstance.Valid {
			userJwtClaim, err := mapper.MapJwtClaimIntoUserClaim(claims)
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
			ginContext.Set("claims", userJwtClaim)
		}

		return nil
	})
	return tokenString
}

func (userService *ServiceImpl) Update(ginContext *gin.Context, updateUserRequest *model.UpdateUserRequest) {
	valErr := userService.validatorService.ValidateStruct(updateUserRequest)
	userService.validatorService.ParseValidationError(valErr, *updateUserRequest)
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		user, err := userService.userRepository.FindById(gormTransaction, updateUserRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdateUserRequestIntoUserEntity(updateUserRequest, user)
		err = userService.userRepository.Update(gormTransaction, user)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (userService *ServiceImpl) Delete(ginContext *gin.Context, deleteUserRequest *model.DeleteUserRequest) {
	valErr := userService.validatorService.ValidateStruct(deleteUserRequest)
	userService.validatorService.ParseValidationError(valErr, *deleteUserRequest)
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := userService.userRepository.Delete(gormTransaction, deleteUserRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
