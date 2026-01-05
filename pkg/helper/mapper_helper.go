package helper

import (
	"fmt"
	"go-intconnect-api/pkg/exception"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
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
		if str, ok := data.(string); ok {
			// Trim untuk keamanan input
			str = strings.TrimSpace(str)

			if str == "" {
				return time.Time{}, nil
			}

			// Daftar layout yang diizinkan
			layouts := []string{
				"2006-01-02",       // YYYY-MM-DD
				"2006-01-02 15:04", // YYYY-MM-DD HH:mm
			}

			var parsed time.Time
			var err error

			for _, layout := range layouts {
				parsed, err = time.Parse(layout, str)
				if err == nil {
					return parsed, nil
				}
			}

			return nil, fmt.Errorf("invalid date format: %s", str)
		}
	}
	return data, nil
}

func DecodeFromSource[S any, T any](sourceMapping S, targetMapping T) T {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook:       StringIntoTypeHookFunc,
		Result:           &targetMapping,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	err = decoder.Decode(sourceMapping)
	return targetMapping
}

func MapEntitiesIntoResponses[S any, R any](entityObjects []S) []R {
	var responseObjects []R
	for _, entityObject := range entityObjects {
		responseObjects = append(responseObjects, MapEntityIntoResponse[S, R](entityObject, nil))
	}
	return responseObjects
}

func MapEntitiesIntoResponsesWithFunc[S any, R any](
	entityObjects []S,
	renderPayloads ...func(S, R),
) []R {
	var responseObjects []R
	for _, entityObject := range entityObjects {
		responseObjects = append(
			responseObjects,
			MapEntityIntoResponse[S, R](entityObject, renderPayloads...),
		)
	}
	return responseObjects
}

func MapEntitiesIntoResponsesWithIgnoredFieldsWithFunc[S any, R any](
	entityObjects []S,
	ignoredFields []string,
	renderPayloads ...func(S, R),
) []R {
	var responseObjects []R
	for _, entityObject := range entityObjects {
		responseObjects = append(
			responseObjects,
			MapEntityIntoResponseWithIgnoredFields[S, R](entityObject, ignoredFields, renderPayloads...),
		)
	}
	return responseObjects
}

func MapEntityIntoResponseWithIgnoredFields[S any, R any](
	entityObject S,
	ignoredFields []string,
	renderPayloads ...func(S, R),
) R {
	var responseObject R
	if reflect.TypeOf(responseObject).Kind() == reflect.Ptr {
		responseObject = reflect.New(reflect.TypeOf(responseObject).Elem()).Interface().(R)
	}
	decoded := DecodeFromSource[S, R](entityObject, responseObject)
	responseObject = decoded

	// Remove fields if necessary
	if len(ignoredFields) > 0 {
		removeIgnoredFields(responseObject, ignoredFields)
	}
	if len(renderPayloads) > 0 {
		for _, renderPayload := range renderPayloads {
			if renderPayload == nil { // <-- cegah panic
				continue
			}
			renderPayload(entityObject, responseObject)
		}
	}
	return responseObject
}

func MapEntityIntoResponse[S any, R any](
	entityObject S,
	renderPayloads ...func(S, R),
) R {
	var responseObject R
	if reflect.TypeOf(responseObject).Kind() == reflect.Ptr {
		responseObject = reflect.New(reflect.TypeOf(responseObject).Elem()).Interface().(R)
	}
	decoded := DecodeFromSource[S, R](entityObject, responseObject)
	responseObject = decoded

	if len(renderPayloads) > 0 {
		for _, renderPayload := range renderPayloads {
			if renderPayload == nil {
				continue
			}
			renderPayload(entityObject, responseObject)
		}
	}
	return responseObject
}

func MapCreateRequestIntoEntity[S any, R any](createRequest *S) *R {
	var entityObject R
	DecodeFromSource[*S, *R](createRequest, &entityObject)
	return &entityObject
}

func MapUpdateRequestIntoEntity[S any, R any](updateRequest S, existingEntity *R) {
	existingEntity = DecodeFromSource[S, *R](updateRequest, existingEntity)
}

func ParsingHashMapIntoStruct[R any](sourceHashMap map[string]interface{}, rawStruct R) *R {
	parsedRawStruct := DecodeFromSource[map[string]interface{}, *R](sourceHashMap, &rawStruct)
	return parsedRawStruct
}

func removeIgnoredFields[R any](obj R, ignored []string) {
	reflectVal := reflect.ValueOf(obj)
	if reflectVal.Kind() == reflect.Ptr {
		reflectVal = reflectVal.Elem()
	}

	for _, field := range ignored {
		fieldName := reflectVal.FieldByName(field)
		if fieldName.IsValid() && fieldName.CanSet() {
			zero := reflect.Zero(fieldName.Type())
			fieldName.Set(zero)
		}
	}
}
