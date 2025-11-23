package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
	"github.com/golang-jwt/jwt/v5"
)

func MapJwtClaimIntoUserClaim(jwtClaim jwt.MapClaims) (*model.JwtClaimRequest, error) {
	var userClaim model.JwtClaimRequest
	err := mapstructure.Decode(jwtClaim, &userClaim)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &userClaim, nil
}

func MapUserEntityIntoUserResponse(userEntity *entity.User) *model.UserResponse {
	var userResponse model.UserResponse
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: helper.StringIntoTypeHookFunc,
		Result:     &userResponse,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))

	err = decoder.Decode(userEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &userResponse
}

func MapUserEntitiesIntoUserResponses(userEntities []entity.User) []*model.UserResponse {
	var userResponses []*model.UserResponse
	for _, userEntity := range userEntities {
		userResponses = append(userResponses, MapUserEntityIntoUserResponse(&userEntity))
	}
	return userResponses
}

func MapCreateUserRequestIntoUserEntity(createUserRequest *model.CreateUserRequest) *entity.User {
	var userEntity entity.User
	err := mapstructure.Decode(createUserRequest, &userEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &userEntity
}

func MapUpdateUserRequestIntoUserEntity(updateUserRequest *model.UpdateUserRequest, userEntity *entity.User) {
	//helper.DecoderConfigMapper(updateUserRequest, &userEntity)
}
