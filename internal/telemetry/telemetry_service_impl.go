package telemetry

import (
	"bytes"
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	reportDocumentTemplate "go-intconnect-api/internal/report_document_template"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
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
		var mapOfParameter = make(map[uint64]*entity.Parameter)
		for _, parameterEntity := range reportDocumentTemplateEntity.Parameters {
			searchedParameterIds = append(searchedParameterIds, parameterEntity.Id)
			mapOfParameterMachine[parameterEntity.Id] = &parameterEntity.MqttTopic.Machine
			mapOfParameter[parameterEntity.Id] = parameterEntity
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
					Timestamp:     telemetryEntity.Bucket,
					MachineId:     mapOfParameterMachine[telemetryEntity.ParameterId].Id,
					MachineName:   mapOfParameterMachine[telemetryEntity.ParameterId].Name,
					MachineCode:   mapOfParameterMachine[telemetryEntity.ParameterId].Code,
					ParameterName: mapOfParameter[telemetryEntity.ParameterId].Name,
					Value:         telemetryEntity.LastValue,
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
func (telemetryService *ServiceImpl) GenerateXLSX(
	ctx *gin.Context,
	req *model.TelemetryReportFilterRequest,
) (bytes.Buffer, error) {

	telemetryGroupeds := telemetryService.GenerateReport(ctx, req)

	f := excelize.NewFile()
	sheetName := "Report"
	f.SetSheetName("Sheet1", sheetName)

	// =====================================================
	// === COPY LAYOUT DARI FILE CONTOH (Book1.xlsx)
	// =====================================================

	// Column width
	_ = f.SetColWidth(sheetName, "A", "A", 24.83203125)
	_ = f.SetColWidth(sheetName, "B", "B", 14)
	_ = f.SetColWidth(sheetName, "C", "C", 13)
	_ = f.SetColWidth(sheetName, "D", "D", 21.33203125)
	_ = f.SetColWidth(sheetName, "F", "F", 26)
	_ = f.SetColWidth(sheetName, "I", "I", 11.1640625)

	// Row height
	_ = f.SetRowHeight(sheetName, 3, 34)

	// Merged cells
	_ = f.MergeCell(sheetName, "A1", "A3")
	_ = f.MergeCell(sheetName, "B2", "D2")
	_ = f.MergeCell(sheetName, "B3", "D3")

	// Header values
	f.SetCellValue(sheetName, "B2", "PT Kalbe Morinaga Indonesia")
	f.SetCellValue(
		sheetName,
		"B3",
		"Jl. Raya Kawasan Industri Indotaisei, Sektor 1A, Blok Q1, Kalihurip, Cikampek, Karawang, West Java 41373",
	)

	// (isi dinamis dari filter kamu)
	f.SetCellValue(
		sheetName,
		"A4",
		fmt.Sprintf(
			"Periode: %s - %s",
			req.StartDate,
			req.EndDate,
		),
	)

	f.SetCellValue(
		sheetName,
		"A5",
		"Generated By: Administrator",
	)

	// Style title & address (supaya sama tampilannya)
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})

	addressStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "center",
		},
	})

	f.SetCellStyle(sheetName, "B2", "D2", titleStyle)
	f.SetCellStyle(sheetName, "B3", "D3", addressStyle)

	// =====================================================
	// === LOGO (AREA A1:A3)
	// =====================================================

	logoPath := "assets/kalbeNutritionalLogo.png"

	if _, err := os.Stat(logoPath); err == nil {

		// hitung scale agar mendekati ukuran inch yang kamu mau
		file, err := os.Open(logoPath)
		if err == nil {
			defer file.Close()

			cfg, _, err := image.DecodeConfig(file)
			if err == nil {

				// target dari kamu
				targetWidthPx := 1.63 * 96
				targetHeightPx := 0.79 * 96

				scaleX := targetWidthPx / float64(cfg.Width)
				scaleY := targetHeightPx / float64(cfg.Height)

				_ = f.AddPicture(
					sheetName,
					"A1",
					logoPath,
					&excelize.GraphicOptions{
						ScaleX:      scaleX,
						ScaleY:      scaleY,
						Positioning: "oneCell",
					},
				)
			}
		}
	}

	// =====================================================
	// === TABLE HEADER
	// =====================================================

	startRow := 7

	headers := []string{
		"Timestamp",
		"Machine Name",
		"Machine Code",
		"Parameter",
		"Value",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, startRow)
		f.SetCellValue(sheetName, cell, h)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	f.SetCellStyle(sheetName, "A7", "E7", headerStyle)

	// =====================================================
	// === CONTENT (flatten dari TelemetryGrouped)
	// =====================================================

	rowIdx := startRow + 1

	for _, grouped := range telemetryGroupeds {

		for _, entry := range grouped.Entries {

			f.SetCellValue(
				sheetName,
				fmt.Sprintf("A%d", rowIdx),
				grouped.Timestamp.Format(time.RFC3339),
			)

			f.SetCellValue(
				sheetName,
				fmt.Sprintf("B%d", rowIdx),
				entry.MachineName,
			)

			f.SetCellValue(
				sheetName,
				fmt.Sprintf("C%d", rowIdx),
				entry.MachineCode,
			)

			f.SetCellValue(
				sheetName,
				fmt.Sprintf("D%d", rowIdx),
				entry.ParameterName,
			)

			if entry.Value != nil {
				f.SetCellValue(
					sheetName,
					fmt.Sprintf("E%d", rowIdx),
					*entry.Value,
				)
			}

			rowIdx++
		}
	}

	// =====================================================
	// === OUTPUT
	// =====================================================

	var buf bytes.Buffer
	err := f.Write(&buf)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.ErrInternalServerError))
	return buf, nil
}

func (telemetryService *ServiceImpl) IntervalReport(ginContext *gin.Context, telemetryIntervalFilterRequest *model.TelemetryIntervalFilterRequest) *model.TelemetryIntervalValues {
	var telemetryIntervalValues *model.TelemetryIntervalValues

	err := telemetryService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {

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

		// 6. Query data dari database
		telemetryEntities, err := telemetryService.telemetryRepository.FindAllInterval(
			gormTransaction,
			telemetryIntervalFilterRequest.ParameterIds,
			intervalVal,
			startDate,
			endDate,
		)
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
