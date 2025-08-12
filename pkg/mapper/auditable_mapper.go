package mapper

import (
	"github.com/go-viper/mapstructure/v2"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
)

func AuditableEntityIntoEntityResponse(auditableEntity *entity.Auditable) *model.AuditableResponse {
	var auditableResponse model.AuditableResponse
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: helper.StringIntoTypeHookFunc,
		Result:     &auditableResponse,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))

	err = decoder.Decode(auditableEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &auditableResponse

}
