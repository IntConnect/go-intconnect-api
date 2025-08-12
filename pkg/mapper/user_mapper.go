package mapper

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/golang-jwt/jwt/v5"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
)

func MapJwtClaimIntoUserClaim(jwtClaim jwt.MapClaims) (*model.JwtClaimDto, error) {
	var userClaim model.JwtClaimDto
	err := mapstructure.Decode(jwtClaim, &userClaim)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &userClaim, nil
}
