package main

import (
	"context"
	"fmt"
	"go-intconnect-api/cmd/injector"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// --- Invoker ---
func Run(fxLifecycle fx.Lifecycle, ginEngine *gin.Engine) {
	fxLifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting server...")
			go func() {
				if err := ginEngine.Run(":8080"); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping server...")
			return nil
		},
	})
}

func main() {

	fxContainer := fx.New(
		// Provider
		injector.CoreModule,
		injector.ApplicationRoutesModule,
		injector.UserModule,
		injector.ValidatorModule,
		injector.NodeModule,
		injector.RoleModule,
		injector.PipelineModule,
		injector.PipelineEdgeModule,
		injector.PipelineNodeModule,
		injector.PipelineConfigurationModule,
		injector.DatabaseConnectionModule,
		injector.FacilityModule,
		injector.PermissionModule,
		injector.LoggerModule,
		injector.MqttBrokerModule,
		injector.MachineModule,
		// Invoker
		fx.Invoke(Run),
	)

	fxContainer.Run()
}
