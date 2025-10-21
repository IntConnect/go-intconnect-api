package injector

import (
	"go-intconnect-api/internal/node"
	"go-intconnect-api/internal/pipeline"
	pipelineEdge "go-intconnect-api/internal/pipeline_edge"
	pipelineNode "go-intconnect-api/internal/pipeline_node"
	protocolConfiguration "go-intconnect-api/internal/protocol_configuration"
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

var AuthenticationRoutesModule = fx.Module("authenticationRoutes",
	fx.Provide(routes.NewAuthenticationRoutes),
	fx.Invoke(func(authenticationRoutes *routes.AuthenticationRoutes) {
		authenticationRoutes.Setup()
	}),
)

var UserModule = fx.Module("userFeature",
	fx.Provide(fx.Annotate(user.NewRepository, fx.As(new(user.Repository)))),
	fx.Provide(fx.Annotate(user.NewService, fx.As(new(user.Service)))),
	fx.Provide(fx.Annotate(user.NewHandler, fx.As(new(user.Controller)))),
)

var NodeModule = fx.Module("nodeFeature",
	fx.Provide(fx.Annotate(node.NewRepository, fx.As(new(node.Repository)))),
	fx.Provide(fx.Annotate(node.NewService, fx.As(new(node.Service)))),
	fx.Provide(fx.Annotate(node.NewHandler, fx.As(new(node.Controller)))),
)

var ValidatorModule = fx.Module("validatorFeature",
	fx.Provide(fx.Annotate(validator.NewService, fx.As(new(validator.Service)))),
)

var PipelineModule = fx.Module("pipelineFeature",
	fx.Provide(fx.Annotate(pipeline.NewRepository, fx.As(new(pipeline.Repository)))),
	fx.Provide(fx.Annotate(pipeline.NewService, fx.As(new(pipeline.Service)))),
	fx.Provide(fx.Annotate(pipeline.NewHandler, fx.As(new(pipeline.Controller)))),
)

var PipelineConfigurationModule = fx.Module("pipelineConfigurationFeature",
	fx.Provide(fx.Annotate(protocolConfiguration.NewRepository, fx.As(new(protocolConfiguration.Repository)))),
	fx.Provide(fx.Annotate(protocolConfiguration.NewService, fx.As(new(protocolConfiguration.Service)))),
	fx.Provide(fx.Annotate(protocolConfiguration.NewHandler, fx.As(new(protocolConfiguration.Controller)))),
)

var PipelineNodeModule = fx.Module("pipelineNodeFeature",
	fx.Provide(fx.Annotate(pipelineNode.NewRepository, fx.As(new(pipelineNode.Repository)))),
)

var PipelineEdgeModule = fx.Module("pipelineFeature",
	fx.Provide(fx.Annotate(pipelineEdge.NewRepository, fx.As(new(pipelineEdge.Repository)))),
)
