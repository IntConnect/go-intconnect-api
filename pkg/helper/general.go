package helper

import (
	"fmt"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"gorm.io/gorm"
)

func CheckErrorOperation(indicatedError error, applicationError *exception.ApplicationError) bool {
	if indicatedError != nil {
		panic(applicationError)
		return true
	}
	return false
}

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

func DecoderConfigMapper[S any, R any](sourceMapper S, resultMapper *R) {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: StringIntoTypeHookFunc,
		Result:     resultMapper,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	err = decoder.Decode(sourceMapper)
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
}

func RandomStringGenerator(length int) (code string) {
	var randomizer = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	c := make([]string, length)
	for i := range c {
		numOrAlpha := rand.Intn(2)
		if numOrAlpha == 0 {
			c[i] = strconv.Itoa(randomizer.Intn(10))
		} else {
			c[i] = string(letters[randomizer.Intn(len(letters))])
		}

		code = strings.Join(c, "")
	}
	return
}

func CheckPointerWrapper[T any](targetChecking *T, renderPayload func()) {
	if targetChecking != nil {
		renderPayload()
	}
}

func DecodeFromSource[T any](sourceMapping interface{}, targetMapping T) T {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: StringIntoTypeHookFunc,
		Result:     &targetMapping,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	err = decoder.Decode(sourceMapping)
	return targetMapping

}

func DebugPrintArray(arrayOfPointer []*model.PipelineNodeResponse) {
	for _, pointerElement := range arrayOfPointer {
		fmt.Println(pointerElement.Config.NodeTempId)
		fmt.Println(pointerElement.Config.ProtocolConfigurationId)
	}
}
