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
	"net/http"
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
			mapOfParameterMachine[parameterEntity.Id] = &parameterEntity.MqttTopic.Machine
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
func (telemetryService *ServiceImpl) IntervalReport(ginContext *gin.Context, telemetryIntervalFilterRequest *model.TelemetryIntervalFilterRequest) *model.TelemetryIntervalValues {
	var telemetryIntervalValues *model.TelemetryIntervalValues

	err := telemetryService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		fmt.Println(*telemetryIntervalFilterRequest)

		// 1. Parse ISO timestamp
		baseTime, err := time.Parse(time.RFC3339, telemetryIntervalFilterRequest.Timestamp)
		helper.CheckErrorOperation(err, exception.NewApplicationError(
			http.StatusBadRequest,
			exception.ErrBadRequest,
		))

		// 2. Parse starting hour
		hourLayout := "15:04:05"
		parsedHour, err := time.Parse(hourLayout, telemetryIntervalFilterRequest.StartingHour)
		helper.CheckErrorOperation(err, exception.NewApplicationError(
			http.StatusBadRequest,
			exception.ErrBadRequest,
		))

		// 3. Load timezone WIB (Asia/Jakarta)
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			// Fallback to FixedZone if LoadLocation fails
			loc = time.FixedZone("WIB", 7*3600) // UTC+7
		}

		// 4. Gabungkan date + hour menggunakan timezone WIB
		startDate := time.Date(
			baseTime.Year(),
			baseTime.Month(),
			baseTime.Day(),
			parsedHour.Hour(),
			parsedHour.Minute(),
			0,
			0,
			loc, // Gunakan WIB timezone
		)

		// 5. End date + 24 jam
		endDate := startDate.Add(24 * time.Hour)

		intervalVal := fmt.Sprintf("%d hours", telemetryIntervalFilterRequest.Interval)
		fmt.Println(intervalVal, startDate, endDate, telemetryIntervalFilterRequest.ParameterIds)

		// 6. Query data dari database
		telemetryEntities, err := telemetryService.telemetryRepository.FindAllInterval(
			gormTransaction,
			telemetryIntervalFilterRequest.ParameterIds,
			intervalVal,
			startDate,
			endDate,
		)
		helper.DebugArrPointer(telemetryEntities)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		// 7. Init response
		telemetryIntervalValues = &model.TelemetryIntervalValues{
			Meta: model.TelemetryMeta{
				Date:         baseTime.Format("2006-01-02"),
				Interval:     telemetryIntervalFilterRequest.Interval,
				Timezone:     baseTime.Format("-07:00"),
				StartingHour: telemetryIntervalFilterRequest.StartingHour[:5], // "HH:mm"
			},
			Data: map[string]map[uint64]*float64{},
		}

		// 8. Inisialisasi map dengan timezone lokal (WIB)
		for t := startDate; t.Before(endDate); t = t.Add(time.Duration(telemetryIntervalFilterRequest.Interval) * time.Hour) {
			timeKey := fmt.Sprintf("%02d:00", t.Hour())
			fmt.Printf("Initializing timeKey: %s (from %s)\n", timeKey, t.String())

			telemetryIntervalValues.Data[timeKey] = map[uint64]*float64{}
			for _, paramID := range telemetryIntervalFilterRequest.ParameterIds {
				telemetryIntervalValues.Data[timeKey][paramID] = nil
			}
		}

		// 9. Mapping data dari database ke response
		for _, telemetryEntity := range telemetryEntities {
			// Konversi bucket time ke timezone lokal (WIB)
			localBucket := telemetryEntity.Bucket.In(loc)
			timeKey := fmt.Sprintf("%02d:00", localBucket.Hour())

			fmt.Printf("Database timeKey: %s (from %s, UTC: %s)\n",
				timeKey,
				localBucket.String(),
				telemetryEntity.Bucket.UTC().String())

			// Pastikan timeKey ada di map
			if slotMap, ok := telemetryIntervalValues.Data[timeKey]; ok {
				slotMap[telemetryEntity.ParameterId] = telemetryEntity.LastValue
			} else {
				// Optional: log jika ada data di luar range
				fmt.Printf("Warning: timeKey %s not found in initialized map\n", timeKey)
			}
		}

		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))

	return telemetryIntervalValues
}
