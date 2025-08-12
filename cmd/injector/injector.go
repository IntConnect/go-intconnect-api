package injector

import (
	"go-intconnect-api/internal/user"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/routes"
	"go.uber.org/fx"
)

var ProtectedRoutesModule = fx.Module("protectedRoutes",
	fx.Provide(routes.NewProtectedRoutes),
	fx.Invoke(func(protectedRoutes *routes.ProtectedRoutes) {
		protectedRoutes.Setup()
	}),
)

var UserModule = fx.Module("userFeature",
	fx.Provide(fx.Annotate(user.NewRepository, fx.As(new(user.Repository)))),
	fx.Provide(fx.Annotate(user.NewService, fx.As(new(user.Service)))),
	fx.Provide(fx.Annotate(user.NewHandler, fx.As(new(user.Controller)))),
)

var ValidatorModule = fx.Module("validatorFeature",
	fx.Provide(fx.Annotate(validator.NewService, fx.As(new(validator.Service)))),
)
