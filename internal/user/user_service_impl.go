package user

import (
	"fmt"
	"go-intconnect-api/configs"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/role"
	"go-intconnect-api/internal/storage"
	"go-intconnect-api/internal/trait"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	userRepository      Repository
	roleService         role.Service
	auditLogService     auditLog.Service
	validatorService    validator.Service
	dbConnection        *gorm.DB
	viperConfig         *viper.Viper
	redisInstance       *configs.RedisInstance
	localStorageService *storage.Manager
}

func NewService(userRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	auditLogService auditLog.Service,
	redisInstance *configs.RedisInstance,
	roleService role.Service,
	localStorageService *storage.Manager,
) *ServiceImpl {
	return &ServiceImpl{
		userRepository:      userRepository,
		validatorService:    validatorService,
		dbConnection:        dbConnection,
		viperConfig:         viperConfig,
		auditLogService:     auditLogService,
		redisInstance:       redisInstance,
		roleService:         roleService,
		localStorageService: localStorageService,
	}
}

func (userService *ServiceImpl) FindAll() []*model.UserResponse {
	var allUser []*model.UserResponse
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userResponse, err := userService.userRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allUser = helper.MapEntitiesIntoResponsesWithFunc[*entity.User, *model.UserResponse](userResponse, mapper.FuncMapAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allUser
}

func (userService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.UserResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var userResponses []*model.UserResponse
	var totalItems int64

	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntities, total, err := userService.userRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		userResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.User, *model.UserResponse](
			userEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Users fetched successfully",
		userResponses,
		paginationReq,
		totalItems,
	)
}

func (userService *ServiceImpl) FindById(ginContext *gin.Context, userId uint64) *model.UserResponse {
	var userResponse *model.UserResponse
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntity, err := userService.userRepository.FindById(gormTransaction, userId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		userResponse = helper.MapEntityIntoResponse[*entity.User, *model.UserResponse](userEntity,
			mapper.FuncMapAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return userResponse
}

func (userService *ServiceImpl) FindSelf(ginContext *gin.Context) *model.UserResponse {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var userResponse *model.UserResponse
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntity, err := userService.userRepository.FindById(gormTransaction, userJwtClaims.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		roleResponse := userService.roleService.FindById(ginContext, userEntity.RoleId)
		userResponse = helper.MapEntityIntoResponse[*entity.User, *model.UserResponse](userEntity,
			mapper.FuncMapAuditable)
		userResponse.RoleResponse = roleResponse
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return userResponse
}

// Create - Membuat user baru
func (userService *ServiceImpl) Create(ginContext *gin.Context, createUserRequest *model.CreateUserRequest) *model.PaginatedResponse[*model.UserResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginationResp *model.PaginatedResponse[*model.UserResponse]
	valErr := userService.validatorService.ValidateStruct(createUserRequest)
	userService.validatorService.ParseValidationError(valErr, *createUserRequest)
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntity := helper.MapCreateRequestIntoEntity[model.CreateUserRequest, entity.User](createUserRequest)
		userEntity.Status = trait.UserStatusActive
		userEntity.Auditable = entity.NewAuditable(userJwtClaims.Name)
		err := userService.userRepository.Create(gormTransaction, userEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := userService.auditLogService.Build(
			nil,
			userEntity,
			nil,
			"",
		)

		err = userService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_CREATE,
				model.AUDIT_LOG_FEATURE_USER,
				auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = userService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (userService *ServiceImpl) HandleLogin(ginContext *gin.Context, loginUserRequest *model.LoginUserRequest) string {
	ipAddress, _ := ginContext.Get("ipAddress")
	err := userService.validatorService.ValidateStruct(loginUserRequest)
	userService.validatorService.ParseValidationError(err, *loginUserRequest)
	var tokenString string
	err = userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntity, err := userService.userRepository.FindByIdentifier(gormTransaction, loginUserRequest.UserIdentifier)

		if err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(loginUserRequest.Password)); err != nil {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, "User credentials invalid"))
		}
		jwtHour := userService.viperConfig.GetInt64("JWT_HOUR")
		if jwtHour == 0 {
			jwtHour = 72
		}
		tokenInstance := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":          userEntity.Id,
			"email":       userEntity.Email,
			"username":    userEntity.Username,
			"name":        userEntity.Name,
			"role_id":     userEntity.RoleId,
			"role_name":   userEntity.Role.Name,
			"avatar_path": userEntity.AvatarPath,
			"exp":         time.Now().Add(time.Hour * time.Duration(jwtHour)).Unix(),
		})
		tokenString, err = tokenInstance.SignedString([]byte(userService.viperConfig.GetString("JWT_SECRET")))
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError))
		if claims, ok := tokenInstance.Claims.(jwt.MapClaims); ok && tokenInstance.Valid {
			userJwtClaim := helper.MapCreateRequestIntoEntity[jwt.MapClaims, *model.JwtClaimRequest](&claims)
			ginContext.Set("claims", userJwtClaim)
		}
		redisKey := fmt.Sprintf("auth:token:%d", userEntity.Id)
		err = userService.redisInstance.RedisClient.Set(
			ginContext.Request.Context(),
			redisKey,
			tokenString,
			time.Hour*time.Duration(jwtHour),
		).Err()
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError))
		userService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      userEntity.Id,
			Action:      model.AUDIT_LOG_LOGIN,
			Feature:     model.AUDIT_LOG_FEATURE_USER,
			Description: "",
			Before:      nil,
			After:       nil,
			IpAddress:   ipAddress.(string),
		})
		return nil
	})
	return tokenString
}
func (userService *ServiceImpl) HandleLogout(ginContext *gin.Context) {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	redisKey := fmt.Sprintf("auth:token:%d", userJwtClaims.Id)
	deletedKey, err := userService.redisInstance.RedisClient.Del(ginContext.Request.Context(), redisKey).Result()
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError))
	if deletedKey <= 0 {
		exception.ThrowApplicationError(exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError))
	}
}
func (userService *ServiceImpl) Update(ginContext *gin.Context, updateUserRequest *model.UpdateUserRequest) *model.PaginatedResponse[*model.UserResponse] {
	var paginationResp *model.PaginatedResponse[*model.UserResponse]
	valErr := userService.validatorService.ValidateStruct(updateUserRequest)
	userService.validatorService.ParseValidationError(valErr, *updateUserRequest)
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntity, err := userService.userRepository.FindById(gormTransaction, updateUserRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		pastUserEntity := *userEntity
		helper.MapUpdateRequestIntoEntity(updateUserRequest, userEntity)
		err = userService.userRepository.Update(gormTransaction, userEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := userService.auditLogService.Build(
			&pastUserEntity,
			userEntity,
			nil,
			"",
		)

		err = userService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_UPDATE,
				model.AUDIT_LOG_FEATURE_USER,
				auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = userService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (userService *ServiceImpl) UpdateProfile(ginContext *gin.Context, updateUserProfileRequest *model.UpdateUserProfileRequest) string {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	updateUserProfileRequest.Id = userJwtClaims.Id
	valErr := userService.validatorService.ValidateStruct(updateUserProfileRequest)
	userService.validatorService.ParseValidationError(valErr, *updateUserProfileRequest)
	var tokenString string
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var err error
		userEntity, err := userService.userRepository.FindById(gormTransaction, updateUserProfileRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateUserProfileRequest, userEntity)
		if updateUserProfileRequest.Password != nil {
			userEntity.Password, _ = helper.HashPassword(*updateUserProfileRequest.Password)
		}
		if updateUserProfileRequest.Avatar != nil {
			newPath, err := userService.localStorageService.Disk().Put(updateUserProfileRequest.Avatar, fmt.Sprintf("users/profiles/%d-%s", time.Now().UnixNano(), updateUserProfileRequest.Avatar.Filename))
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrSavingResources))
			userEntity.AvatarPath = newPath
		}

		err = userService.userRepository.Update(gormTransaction, userEntity)
		tokenString = userService.HandleLogin(ginContext, &model.LoginUserRequest{
			UserIdentifier: userEntity.Email,
			Password:       updateUserProfileRequest.CurrentPassword,
		})
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return tokenString
}

func (userService *ServiceImpl) Delete(ginContext *gin.Context, deleteUserRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.UserResponse] {
	var paginationResp *model.PaginatedResponse[*model.UserResponse]
	valErr := userService.validatorService.ValidateStruct(deleteUserRequest)
	userService.validatorService.ParseValidationError(valErr, *deleteUserRequest)
	err := userService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userEntity, err := userService.userRepository.FindById(gormTransaction, deleteUserRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = userService.userRepository.Delete(gormTransaction, deleteUserRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		auditPayload := userService.auditLogService.Build(
			userEntity,
			nil,
			nil,
			deleteUserRequest.Reason,
		)

		err = userService.auditLogService.
			Record(ginContext,
				model.AUDIT_LOG_DELETE,
				model.AUDIT_LOG_FEATURE_USER,
				auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = userService.FindAllPagination(&paginationRequest)
	return paginationResp
}
