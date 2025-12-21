package telemetry

import (
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	reportDocumentTemplate "go-intconnect-api/internal/report_document_template"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	telemetryRepository    Repository
	reportDocumentTemplate reportDocumentTemplate.Repository
	auditLogService        auditLog.Service
	validatorService       validator.Service
	dbConnection           *gorm.DB
	viperConfig            *viper.Viper
}

func NewService(telemetryRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	auditLogService auditLog.Service,
	reportDocumentTemplate reportDocumentTemplate.Repository,

) *ServiceImpl {
	return &ServiceImpl{
		telemetryRepository:    telemetryRepository,
		validatorService:       validatorService,
		dbConnection:           dbConnection,
		viperConfig:            viperConfig,
		auditLogService:        auditLogService,
		reportDocumentTemplate: reportDocumentTemplate,
	}
}

func (telemetryService *ServiceImpl) GenerateReport(ginContext *gin.Context, telemetryReportFilterRequest *model.TelemetryReportFilterRequest) []*model.TelemetryGrouped {
	var telemetriesGrouped []*model.TelemetryGrouped
	err := telemetryService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		reportDocumentTemplateEntity, err := telemetryService.reportDocumentTemplate.FindById(gormTransaction, telemetryReportFilterRequest.ReportDocumentTemplateId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		var searchedParameterIds []uint64
		var mapOfParameterMachine = make(map[uint64]*entity.Machine)
		for _, parameterEntity := range reportDocumentTemplateEntity.Parameters {
			searchedParameterIds = append(searchedParameterIds, parameterEntity.Id)
			mapOfParameterMachine[parameterEntity.Id] = parameterEntity.MqttTopic.Machine
		}
		dateTimeLayout := "2006-01-02 15:04"
		startDate, _ := time.Parse(dateTimeLayout, telemetryReportFilterRequest.StartDate)
		endDate, _ := time.Parse(dateTimeLayout, telemetryReportFilterRequest.EndDate)
		intervalVal := fmt.Sprintf("%d minutes", telemetryReportFilterRequest.Interval)
		telemetryEntities, err := telemetryService.telemetryRepository.FindAllFilter(gormTransaction, searchedParameterIds, intervalVal, startDate, endDate)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		var mapOfTelemetryReportValue = make(map[time.Time][]*model.TelemetryReportValue)
		for _, telemetryEntity := range telemetryEntities {
			if telemetryReportValueArr, isExists := mapOfTelemetryReportValue[telemetryEntity.Bucket]; isExists {
				telemetryReportValueArr = append(telemetryReportValueArr, &model.TelemetryReportValue{
					Timestamp:   telemetryEntity.Bucket,
					MachineId:   mapOfParameterMachine[telemetryEntity.ParameterId].Id,
					MachineName: mapOfParameterMachine[telemetryEntity.ParameterId].Name,
					MachineCode: mapOfParameterMachine[telemetryEntity.ParameterId].Code,
					ParameterId: telemetryEntity.ParameterId,
					Value:       telemetryEntity.LastValue,
				})
				mapOfTelemetryReportValue[telemetryEntity.Bucket] = telemetryReportValueArr
			} else {
				mapOfTelemetryReportValue[telemetryEntity.Bucket] = []*model.TelemetryReportValue{}
			}
		}
		for timestampBucket, telemetryReportValue := range mapOfTelemetryReportValue {
			telemetriesGrouped = append(telemetriesGrouped, &model.TelemetryGrouped{
				Timestamp: timestampBucket,
				Entries:   telemetryReportValue,
			})
		}

		sort.Slice(telemetriesGrouped, func(i, j int) bool {
			return telemetriesGrouped[i].Timestamp.Before(telemetriesGrouped[j].Timestamp) // Balik tanda >
		})

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return telemetriesGrouped
}
