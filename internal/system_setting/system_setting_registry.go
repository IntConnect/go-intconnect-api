package system_setting

import (
	"go-intconnect-api/pkg/exception"
)

type Registry struct {
	systemSettingHandlers map[string]SystemSettingHandler
}

func NewRegistry() *Registry {
	return &Registry{
		systemSettingHandlers: make(map[string]SystemSettingHandler),
	}
}

func (handlerRegistry *Registry) Register(key string, handler SystemSettingHandler) {
	handlerRegistry.systemSettingHandlers[key] = handler
}

func (handlerRegistry *Registry) Resolve(key string) SystemSettingHandler {
	handler, exists := handlerRegistry.systemSettingHandlers[key]
	if !exists {
		exception.ThrowApplicationError(
			exception.NewApplicationError(400, exception.ErrSystemSettingKeyNotMatch),
		)
	}
	return handler
}
