package main

import (
	"context"
	"fmt"
	"go-intconnect-api/cmd/injector"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func Run(fxLifecycle fx.Lifecycle, ginEngine *gin.Engine, viperConfig *viper.Viper) {
	webPort := viperConfig.GetString("API_PORT")
	fxLifecycle.Append(fx.Hook{

		OnStart: func(ctx context.Context) error {
			logrus.Debug("Starting server...")
			go func() {
				if err := ginEngine.Run(fmt.Sprintf(":%s", webPort)); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Debug("Stopping server...")
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
		injector.MqttBrokerModule,
		injector.MachineModule,
		injector.ParameterModule,
		injector.MachineDocumentModule,
		injector.MqttTopicModule,
		injector.ReportDocumentTemplateModule,
		injector.AuditLogModule,
		injector.SmtpServerModule,
		injector.ModbusServerModule,
		injector.BreakdownModule,
		injector.CheckSheetDocumentModule,
		injector.BreakdownResourceModule,
		injector.SystemSettingModule,
		// Invoker
		fx.Invoke(Run),
	)

	fxContainer.Run()
}
