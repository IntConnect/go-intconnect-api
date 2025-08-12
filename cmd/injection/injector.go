//go:build wireinject
// +build wireinject

package injection

import (
	"github.com/gin-gonic/gin"
	universalTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go-intconnect-api/internal/user"
	"go-intconnect-api/routes"
	"gorm.io/gorm"
)

var userSet = wire.NewSet(
	user.NewRepository,
	wire.Bind(new(user.Repository), new(*user.RepositoryImpl)),
	user.NewService,
	wire.Bind(new(user.Service), new(*user.ServiceImpl)),
	user.NewHandler,
	wire.Bind(new(user.Controller), new(*user.Handler)),
)

// wire.go
func InitializeRoutes(
	ginRouterGroup *gin.RouterGroup,
	dbConnection *gorm.DB,
	validatorInstance *validator.Validate,
	engTranslator universalTranslator.Translator,
	viperConfig *viper.Viper,
) (*routes.ApplicationRoutes, error) {
	wire.Build(
		wire.Struct(new(routes.ApplicationRoutes), "*"),
		userSet,
	)
	return nil, nil
}
