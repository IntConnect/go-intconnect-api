package helper

import (
	"fmt"
	"go-intconnect-api/pkg/exception"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"gorm.io/gorm"
)

func StringIntoTypeHookFunc(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	switch to {
	case reflect.TypeOf(uint64(0)):
		if str, ok := data.(string); ok {
			return strconv.ParseUint(str, 10, 64)
		}
	case reflect.TypeOf(uint32(0)):
		if str, ok := data.(string); ok {
			val, err := strconv.ParseUint(str, 10, 32)
			return uint32(val), err
		}
	case reflect.TypeOf(int(0)):
		if str, ok := data.(string); ok {
			return strconv.Atoi(str)
		}
	case reflect.TypeOf(float32(0)):
		if str, ok := data.(string); ok {
			val, err := strconv.ParseFloat(str, 32)
			return float32(val), err
		}
	case reflect.TypeOf(time.Time{}):
		if rawTime, ok := data.(time.Time); ok {
			return rawTime.String(), nil
		}
	case reflect.TypeOf(gorm.DeletedAt{}):
		if gormDeletedAt, ok := data.(gorm.DeletedAt); ok {
			return gormDeletedAt.Time.String(), nil
		} else {
			return nil, nil
		}
	}
	return data, nil
}

func DecodeFromSource[S any, T any](sourceMapping S, targetMapping T) T {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: StringIntoTypeHookFunc,
		Result:     &targetMapping,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err.Error(), nil))
	err = decoder.Decode(sourceMapping)
	return targetMapping
}

func MapEntitiesIntoResponses[S any, R any](entityObjects []S) []*R {
	var responseObjects []*R
	for _, entityObject := range entityObjects {
		fmt.Println(entityObject)
		//responseObjects = append(responseObjects, MapEntityIntoResponse[S, R](entityObject, nil))
	}
	return responseObjects
}

// 1. Ubah helper function signature
func MapEntitiesIntoResponsesWithFunc[S any, R any](
	entityObjects []S,
	individualRenderPayload func(S, R), // ✅ Ubah dari func(S, *R) *R
) []R { // ✅ Return []R bukan []*R
	var responseObjects []R
	for _, entityObject := range entityObjects {
		responseObjects = append(
			responseObjects,
			MapEntityIntoResponse[S, R](entityObject, individualRenderPayload),
		)
	}
	return responseObjects
}

func MapEntityIntoResponse[S any, R any](
	entityObject S,
	renderPayload func(S, R),
) R {
	var responseObject R

	if reflect.TypeOf(responseObject).Kind() == reflect.Ptr {
		responseObject = reflect.New(reflect.TypeOf(responseObject).Elem()).Interface().(R)
	}

	decoded := DecodeFromSource[S, R](entityObject, responseObject)
	responseObject = decoded

	if renderPayload != nil {
		renderPayload(entityObject, responseObject)
	}

	return responseObject
}
func MapCreateRequestIntoEntity[S any, R any](createRequest *S) *R {
	var entityObject R
	err := mapstructure.Decode(createRequest, &entityObject)
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err.Error(), nil))
	return &entityObject
}

func MapUpdateRequestIntoEntity[S any, R any](updateRequest S, existingEntity *R) {
	existingEntity = DecodeFromSource[S, *R](updateRequest, existingEntity)
}
