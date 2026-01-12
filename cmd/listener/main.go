package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-intconnect-api/configs"
	websocketConfig "go-intconnect-api/configs/websocket"
	"go-intconnect-api/internal/trait"
	"go-intconnect-api/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
	mqttTopic "go-intconnect-api/internal/mqtt_topic"
	"go-intconnect-api/internal/parameter"
	"go-intconnect-api/internal/telemetry"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"

	"go-intconnect-api/cmd/injector"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

const (
	InsertionWorkersAmount = 25
	InsertionQueueSize     = 1000
)

func main() {
	contextWithCancel, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	listenerFluxorInstance := NewListenerFluxor()
	listenerFluxorInstance.StartPeriodicChecker()
	listenerFluxorInstance.StartWorkers(contextWithCancel)
	listenerFluxorInstance.StartTopicHandler(contextWithCancel)
	listenerFluxorInstance.StartSnapshotSaver(contextWithCancel)
	go listenerFluxorInstance.websocketHub.Run()
	go startWebSocketServer(listenerFluxorInstance.websocketHub)

	listenerFluxorInstance.WaitForShutdown()

}

type AbnormalWindow struct {
	Values        []float64
	LastCheckedAt time.Time
}

type AlarmRecoveryWindow struct {
	NormalCount int
}

type ListenerFluxor struct {
	gormDatabase                   *gorm.DB
	mqttBrokersMap                 map[uint64]*entity.MqttBroker
	parametersMap                  map[string]*entity.Parameter
	processedParametersMap         map[uint64]*entity.Parameter
	processedParameterSequencesMap map[uint64][]entity.ProcessedParameterSequence
	insertionChan                  chan []*entity.Telemetry
	waitGroup                      *sync.WaitGroup
	telemetryRepository            telemetry.Repository
	parameterRepository            parameter.Repository
	mqttBrokerRepository           mqttBroker.Repository
	mqttTopicRepository            mqttTopic.Repository
	mqttClient                     mqtt.Client
	rwMutex                        sync.RWMutex
	latestTelemetry                map[string]*entity.Telemetry
	telemetryMutex                 sync.Mutex
	abnormalBuffers                map[uint64]*AbnormalWindow
	recoveryBuffers                map[uint64]*AlarmRecoveryWindow
	alarmMutex                     sync.Mutex
	redisInstance                  *configs.RedisInstance
	websocketHub                   *websocketConfig.Hub
}

func NewListenerFluxor() *ListenerFluxor {
	viperConfig := injector.NewViperConfig()
	databaseCredentials := injector.NewDatabaseCredentials(viperConfig)
	gormDatabase := injector.NewDatabaseConnection(databaseCredentials)
	telemetryRepository := telemetry.NewRepository()
	parameterRepository := parameter.NewRepository()
	mqttBrokerRepository := mqttBroker.NewRepository()
	mqttTopicRepository := mqttTopic.NewRepository()
	redisHostName, redisPassword, redisDatabaseIndex := configs.LoadRedisConfigFromEnvironment(viperConfig)
	redisConfig := configs.NewRedisConfig(redisHostName, redisPassword, redisDatabaseIndex)
	redisInstance, err := configs.InitRedisInstance(redisConfig)
	websocketHub := websocketConfig.NewHub()

	if err != nil {
		logger.WithError(err).Fatal("failed to init redis instance")
	}
	latestTelemetry := make(map[string]*entity.Telemetry)
	listenerFluxor := &ListenerFluxor{
		gormDatabase:                   gormDatabase,
		mqttBrokersMap:                 make(map[uint64]*entity.MqttBroker),
		parametersMap:                  make(map[string]*entity.Parameter),
		processedParametersMap:         make(map[uint64]*entity.Parameter),
		processedParameterSequencesMap: make(map[uint64][]entity.ProcessedParameterSequence),
		insertionChan:                  make(chan []*entity.Telemetry, InsertionQueueSize),
		waitGroup:                      &sync.WaitGroup{},
		telemetryRepository:            telemetryRepository,
		parameterRepository:            parameterRepository,
		mqttBrokerRepository:           mqttBrokerRepository,
		mqttTopicRepository:            mqttTopicRepository,
		latestTelemetry:                latestTelemetry,
		abnormalBuffers:                make(map[uint64]*AbnormalWindow),
		recoveryBuffers:                make(map[uint64]*AlarmRecoveryWindow),
		redisInstance:                  redisInstance,
		websocketHub:                   websocketHub,
	}

	if err := listenerFluxor.loadInitialConfiguration(); err != nil {
		logger.WithError(err).Warn("Failed to load initial configuration")
	}

	if err := listenerFluxor.startMqttConnection(); err != nil {
		logger.WithError(err).Error("Failed to connect to mqtt brokers at startup")
	}

	return listenerFluxor
}

func (listenerFluxor *ListenerFluxor) loadInitialConfiguration() error {
	logger.Info("Loading initial configuration from DB...")
	mqttBrokerEntities, err := listenerFluxor.mqttBrokerRepository.FindAll(listenerFluxor.gormDatabase)
	if err != nil {
		return err
	}
	parameterEntities, err := listenerFluxor.parameterRepository.FindAll(listenerFluxor.gormDatabase, nil)
	if err != nil {
		return err
	}

	// Load processed parameter sequences
	var processedSequences []entity.ProcessedParameterSequence
	if err := listenerFluxor.gormDatabase.Model(&entity.ProcessedParameterSequence{}).Preload("Parameter").Find(&processedSequences).Error; err != nil {
		logger.WithError(err).Warn("Failed to load processed parameter sequences")
	}

	listenerFluxor.rwMutex.Lock()
	defer listenerFluxor.rwMutex.Unlock()

	for _, mqttBrokerEntity := range mqttBrokerEntities {
		listenerFluxor.mqttBrokersMap[mqttBrokerEntity.Id] = mqttBrokerEntity
	}

	for _, parameterEntity := range parameterEntities {
		listenerFluxor.parametersMap[parameterEntity.Code] = parameterEntity

		// Track processed parameters separately
		if parameterEntity.IsProcessed {
			listenerFluxor.processedParametersMap[parameterEntity.Id] = parameterEntity
		}
	}

	// Group sequences by parent parameter ID
	for _, seq := range processedSequences {
		listenerFluxor.processedParameterSequencesMap[seq.ParentParameterId] = append(
			listenerFluxor.processedParameterSequencesMap[seq.ParentParameterId],
			seq,
		)
	}

	logger.Infof("Loaded %d mqtt brokers, %d parameters, %d processed parameters",
		len(listenerFluxor.mqttBrokersMap),
		len(listenerFluxor.parametersMap),
		len(listenerFluxor.processedParametersMap))
	return nil
}

func (listenerFluxor *ListenerFluxor) startMqttConnection() error {
	listenerFluxor.rwMutex.RLock()
	defer listenerFluxor.rwMutex.RUnlock()

	mqttClientOptions := mqtt.NewClientOptions()
	mqttClientOptions.SetClientID(fmt.Sprintf("listenerFluxor-%d", time.Now().UnixNano()))
	mqttClientOptions.AutoReconnect = true
	mqttClientOptions.ConnectTimeout = 5 * time.Second

	added := 0
	for _, mqttBrokerEntity := range listenerFluxor.mqttBrokersMap {
		mqttClientOptions.AddBroker(fmt.Sprintf("tcp://%s:%s", mqttBrokerEntity.HostName, mqttBrokerEntity.MqttPort))
		added++
	}
	if added == 0 {
		logger.Warn("No mqtt brokers configured yet")
		return nil
	}

	mqttClient := mqtt.NewClient(mqttClientOptions)
	token := mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	listenerFluxor.mqttClient = mqttClient
	logger.Info("MQTT client connected")
	return nil
}

func (listenerFluxor *ListenerFluxor) StartPeriodicChecker() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.WithField("panic", r).Error("Recovered from panic in periodic checker")
					}
				}()
				listenerFluxor.CheckConfigurationPeriodically()
			}()
		}
	}()
}

func (listenerFluxor *ListenerFluxor) CheckConfigurationPeriodically() {
	logger.Info("Running periodic configuration check...")

	mqttBrokerEntities, err := listenerFluxor.mqttBrokerRepository.FindAll(listenerFluxor.gormDatabase)
	if err != nil {
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return
	}
	parameterEntities, err := listenerFluxor.parameterRepository.FindAll(listenerFluxor.gormDatabase, nil)
	if err != nil {
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return
	}

	// Reload processed parameter sequences
	var processedSequences []entity.ProcessedParameterSequence
	if err := listenerFluxor.gormDatabase.Preload("Parameter").Find(&processedSequences).Error; err != nil {
		logger.WithError(err).Warn("Failed to reload processed parameter sequences")
	}

	newMqttBrokersMap := make(map[uint64]*entity.MqttBroker)
	newParametersMap := make(map[string]*entity.Parameter)
	newProcessedParametersMap := make(map[uint64]*entity.Parameter)
	newProcessedSequencesMap := make(map[uint64][]entity.ProcessedParameterSequence)

	for _, mqttBrokerEntity := range mqttBrokerEntities {
		newMqttBrokersMap[mqttBrokerEntity.Id] = mqttBrokerEntity
	}

	for _, parameterEntity := range parameterEntities {
		newParametersMap[parameterEntity.Code] = parameterEntity
		if parameterEntity.IsProcessed {
			newProcessedParametersMap[parameterEntity.Id] = parameterEntity
		}
	}

	for _, seq := range processedSequences {
		newProcessedSequencesMap[seq.ParentParameterId] = append(
			newProcessedSequencesMap[seq.ParentParameterId],
			seq,
		)
	}

	listenerFluxor.rwMutex.RLock()
	changed := isMqttBrokerMapChanged(listenerFluxor.mqttBrokersMap, newMqttBrokersMap)
	listenerFluxor.rwMutex.RUnlock()

	if changed {
		logger.Warn("MQTT broker configuration changed. Restarting MQTT broker...")
		listenerFluxor.rwMutex.Lock()
		listenerFluxor.mqttBrokersMap = newMqttBrokersMap
		listenerFluxor.parametersMap = newParametersMap
		listenerFluxor.processedParametersMap = newProcessedParametersMap
		listenerFluxor.processedParameterSequencesMap = newProcessedSequencesMap
		listenerFluxor.rwMutex.Unlock()

		listenerFluxor.RestartMqttBroker()
	} else {
		listenerFluxor.rwMutex.Lock()
		listenerFluxor.parametersMap = newParametersMap
		listenerFluxor.processedParametersMap = newProcessedParametersMap
		listenerFluxor.processedParameterSequencesMap = newProcessedSequencesMap
		listenerFluxor.rwMutex.Unlock()
	}
}

func (listenerFluxor *ListenerFluxor) RestartMqttBroker() {
	if listenerFluxor.mqttClient != nil && listenerFluxor.mqttClient.IsConnected() {
		logger.Info("Disconnecting old mqtt client...")
		listenerFluxor.mqttClient.Disconnect(250)
	}

	if err := listenerFluxor.startMqttConnection(); err != nil {
		logger.WithError(err).Error("MQTT reconnect failed")
		return
	}

	if err := listenerFluxor.resubscribeTopics(context.Background()); err != nil {
		logger.WithError(err).Error("Failed to resubscribe topics after broker restart")
		return
	}

	logger.Info("MQTT broker restarted and reconnected.")
}

func (listenerFluxor *ListenerFluxor) resubscribeTopics(ctx context.Context) error {
	if listenerFluxor.mqttClient == nil || !listenerFluxor.mqttClient.IsConnected() {
		logger.Warn("MQTT client not connected; skipping resubscribe")
		return fmt.Errorf("mqtt not connected")
	}

	mqttTopics, err := listenerFluxor.mqttTopicRepository.FindAll(listenerFluxor.gormDatabase)
	if err != nil {
		return err
	}

	listenerFluxorResp := converterMqttTopicsToListenerResponse(mqttTopics)

	if len(listenerFluxorResp.SubscribeMultiple) > 0 {
		keys := getTopicKeys(listenerFluxorResp.SubscribeMultiple)
		if token := listenerFluxor.mqttClient.Unsubscribe(keys...); token.Wait() && token.Error() != nil {
			logger.WithError(token.Error()).Warn("Error during unsubscribe (may be ok)")
		}
	}

	if token := listenerFluxor.mqttClient.SubscribeMultiple(listenerFluxorResp.SubscribeMultiple, func(client mqtt.Client, message mqtt.Message) {
		listenerFluxor.onMessageReceived(message, listenerFluxorResp.TopicParameter)
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	logger.Infof("Resubscribed %d topics", len(listenerFluxorResp.SubscribeMultiple))
	return nil
}

func (listenerFluxor *ListenerFluxor) StartWorkers(ctx context.Context) {
	for i := 1; i <= InsertionWorkersAmount; i++ {
		listenerFluxor.waitGroup.Add(1)
		go insertionWorker(i, ctx, listenerFluxor.insertionChan, listenerFluxor.telemetryRepository, listenerFluxor.waitGroup, listenerFluxor.gormDatabase)
	}
}

func (listenerFluxor *ListenerFluxor) StartTopicHandler(ctx context.Context) {
	go listenerFluxor.handledTopics(ctx, listenerFluxor.mqttClient, listenerFluxor.mqttTopicRepository, listenerFluxor.insertionChan, listenerFluxor.gormDatabase)
}

func (listenerFluxor *ListenerFluxor) handledTopics(ctxWithCancel context.Context, mqttClient mqtt.Client, mqttTopicRepository mqttTopic.Repository, insertionChan chan<- []*entity.Telemetry, gormDatabase *gorm.DB) {
	if mqttClient == nil {
		logger.Warn("MQTT client is nil in handledTopics; subscribe will fail until connected")
	}

	topicReloadTicker := time.NewTicker(1 * time.Second)
	defer topicReloadTicker.Stop()

	var subscribedTopics *model.MqttTopicListenerResponse = &model.MqttTopicListenerResponse{
		SubscribeMultiple: make(map[string]byte),
		TopicParameter:    make(model.TopicParameter),
	}

	for {
		select {
		case <-ctxWithCancel.Done():
			logger.Info("handledTopics exiting due to context done")
			return
		case <-topicReloadTicker.C:
			mqttTopics, err := mqttTopicRepository.FindAll(gormDatabase)
			if err != nil {
				logger.WithError(err).Error("Error fetching topics")
				continue
			}
			updatedTopics := converterMqttTopicsToListenerResponse(mqttTopics)

			if subscribedTopics != nil && !isTopicMapEqual(subscribedTopics.SubscribeMultiple, updatedTopics.SubscribeMultiple) {
				if subscribedTopics.SubscribeMultiple != nil && len(subscribedTopics.SubscribeMultiple) > 0 && mqttClient != nil {
					if token := mqttClient.Unsubscribe(getTopicKeys(subscribedTopics.SubscribeMultiple)...); token.Wait() && token.Error() != nil {
						logger.WithError(token.Error()).Error("Error unsubscribing from old topics")
					}
				}
				if mqttClient != nil {
					if token := mqttClient.SubscribeMultiple(updatedTopics.SubscribeMultiple, func(mqttClient mqtt.Client, mqttMessage mqtt.Message) {
						listenerFluxor.onMessageReceived(mqttMessage, updatedTopics.TopicParameter)
					}); token.Wait() && token.Error() != nil {
						logger.WithError(token.Error()).Error("Error subscribing to updated topics")
					} else {
						logger.Infof("Subscribed to %d topics", len(updatedTopics.SubscribeMultiple))
					}
				}
				subscribedTopics = updatedTopics
			}
		}
	}
}

func (listenerFluxor *ListenerFluxor) onMessageReceived(mqttMessage mqtt.Message, topicParameter model.TopicParameter) {
	rawMqttPayload := string(mqttMessage.Payload())
	var mqttPayload model.MqttPayload
	err := json.Unmarshal([]byte(rawMqttPayload), &mqttPayload)
	if err != nil {
		logger.Debug("Error unmarshaling mqtt payload")
		logger.Debug(err)
	}
	listenerFluxor.telemetryMutex.Lock()
	defer listenerFluxor.telemetryMutex.Unlock()
	detailMqttTopic, isExists := topicParameter[mqttMessage.Topic()]
	if !isExists {
		return
	}
	var newParameters []*entity.Parameter
	for mqttKey, _ := range mqttPayload.MqttInnerPayload {
		listenerFluxor.rwMutex.RLock()
		_, isExists := listenerFluxor.parametersMap[mqttKey]
		listenerFluxor.rwMutex.RUnlock()
		mqttTopicId := detailMqttTopic["mqtt_topic_id"]
		if !isExists {

			newParameters = append(newParameters, &entity.Parameter{
				MqttTopicId:      &mqttTopicId,
				Name:             mqttKey,
				Code:             mqttKey,
				Unit:             "N/A",
				MinValue:         0,
				MaxValue:         0,
				Description:      "N/A",
				Category:         trait.ParameterCategoryData,
				PositionX:        0,
				PositionY:        0,
				PositionZ:        0,
				RotationX:        0,
				RotationY:        0,
				RotationZ:        0,
				AbnormalDuration: 1,
				IsAutomatic:      true,
				IsDisplay:        false,
				IsWatch:          false,
				IsRunningTime:    false,
				IsFeatured:       false,
				IsProcessed:      false,
				Auditable:        entity.NewAuditable("System"),
			})
		}
	}

	if len(newParameters) > 0 {
		err := listenerFluxor.parameterRepository.CreateBatch(listenerFluxor.gormDatabase, newParameters)
		if err != nil {
			logger.WithError(err).Error("Error creating new parameters")
		}
		parameterEntities, err := listenerFluxor.parameterRepository.FindAll(listenerFluxor.gormDatabase, nil)
		for _, parameterEntity := range parameterEntities {
			listenerFluxor.parametersMap[parameterEntity.Code] = parameterEntity
		}
	}
	for mqttKey, mqttValue := range mqttPayload.MqttInnerPayload {
		listenerFluxor.rwMutex.RLock()
		parameterEntity, isExists := listenerFluxor.parametersMap[mqttKey]
		listenerFluxor.rwMutex.RUnlock()
		if isExists && len(mqttValue) > 0 {
			var parsedMqttValue float64
			switch rawMqttValue := mqttValue[0].(type) {
			case float64:
				parsedMqttValue = mqttValue[0].(float64)
			case bool:
				if rawMqttValue == true {
					parsedMqttValue = 0
				} else {
					parsedMqttValue = 1
				}
			default:
				continue
			}

			listenerFluxor.latestTelemetry[mqttKey] = &entity.Telemetry{
				ParameterId: parameterEntity.Id,
				Value:       parsedMqttValue,
				Timestamp:   mqttPayload.Timestamp.Time,
			}
			if parameterEntity.IsWatch {
				fmt.Println(parameterEntity.Name)
				listenerFluxor.checkAbnormality(parameterEntity, parsedMqttValue)
			}

		} else {
			logger.WithError(err).Info("Parameter not exists")
		}
	}

}

func insertionWorker(id int, contextWithCancel context.Context, insertionChan <-chan []*entity.Telemetry, telemetryRepository telemetry.Repository, waitGroup *sync.WaitGroup, gormDatabase *gorm.DB) {
	defer waitGroup.Done()
	logger.Infof("Insertion worker %d started", id)

	for {
		select {
		case <-contextWithCancel.Done():
			logger.Infof("Worker %d stopping due to context done", id)
			return
		case telemetryEntities, isExists := <-insertionChan:
			if !isExists {
				logger.Infof("Worker %d: Channel closed. Exiting.", id)
				return
			}

			if err := telemetryRepository.CreateBatch(gormDatabase, telemetryEntities); err != nil {
				logger.WithError(err).Errorf("Worker %d: failed to insert telemetry", id)
			}
		}
	}
}

func (listenerFluxor *ListenerFluxor) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Infof("Shutting down...")

	close(listenerFluxor.insertionChan)

	if listenerFluxor.mqttClient != nil && listenerFluxor.mqttClient.IsConnected() {
		listenerFluxor.mqttClient.Disconnect(250)
	}

	listenerFluxor.waitGroup.Wait()

	logger.Info("Shutdown complete.")
}

func isMqttBrokerMapChanged(oldMap, newMap map[uint64]*entity.MqttBroker) bool {
	if len(oldMap) != len(newMap) {
		return true
	}
	for key, newItem := range newMap {
		oldItem, exists := oldMap[key]
		if !exists {
			return true
		}
		if !reflect.DeepEqual(oldItem, newItem) {
			return true
		}
	}
	return false
}

func isTopicMapEqual(a, b map[string]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if vb, ok := b[k]; !ok || vb != v {
			return false
		}
	}
	return true
}

func getTopicKeys(topics map[string]byte) []string {
	keys := make([]string, 0, len(topics))
	for key := range topics {
		keys = append(keys, key)
	}
	return keys
}

func converterMqttTopicsToListenerResponse(mqttTopicEntities []entity.MqttTopic) *model.MqttTopicListenerResponse {
	mqttTopicListenerResponse := &model.MqttTopicListenerResponse{
		SubscribeMultiple: make(map[string]byte),
		TopicParameter:    make(model.TopicParameter),
	}
	for _, mqttTopicEntity := range mqttTopicEntities {
		var QoS byte = 0
		if mqttTopicEntity.QoS > 0 {
			QoS = byte(mqttTopicEntity.QoS)
		}
		mqttTopicListenerResponse.SubscribeMultiple[mqttTopicEntity.Name] = QoS
		mqttTopicListenerResponse.TopicParameter[mqttTopicEntity.Name] = map[string]uint64{
			"mqtt_topic_id": mqttTopicEntity.Id,
			"machine_id":    mqttTopicEntity.MachineId,
		}
	}
	return mqttTopicListenerResponse
}

func (listenerFluxor *ListenerFluxor) StartSnapshotSaver(ctx context.Context) {
	go func() {
		snapshotTicker := time.NewTicker(3 * time.Second)
		defer snapshotTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-snapshotTicker.C:
				listenerFluxor.saveSnapshot()
			}
		}
	}()
}

func (listenerFluxor *ListenerFluxor) saveSnapshot() {
	listenerFluxor.telemetryMutex.Lock()

	// Create snapshot of current telemetry data
	snapshotPayload := make([]*entity.Telemetry, 0, len(listenerFluxor.latestTelemetry))
	for _, latestTelemetry := range listenerFluxor.latestTelemetry {
		snapshotPayload = append(snapshotPayload, latestTelemetry)
	}

	// Calculate processed parameters before saving
	processedTelemetry := listenerFluxor.calculateProcessedParameters()
	snapshotPayload = append(snapshotPayload, processedTelemetry...)

	listenerFluxor.telemetryMutex.Unlock()

	if len(snapshotPayload) == 0 {
		return
	}

	if err := listenerFluxor.telemetryRepository.CreateBatch(listenerFluxor.gormDatabase, snapshotPayload); err != nil {
		logger.WithError(err).Error("Failed to save snapshot telemetry")
	} else {
		logger.Infof("Saved %d telemetry records (snapshot)", len(snapshotPayload))
	}
}

// calculateProcessedParameters computes values for all processed parameters
func (listenerFluxor *ListenerFluxor) calculateProcessedParameters() []*entity.Telemetry {
	listenerFluxor.rwMutex.RLock()
	defer listenerFluxor.rwMutex.RUnlock()

	var processedTelemetry []*entity.Telemetry

	// Create a map for quick lookup by parameter code
	telemetryByCode := make(map[string]float64)
	for parameterCode, latestTelemetry := range listenerFluxor.latestTelemetry {
		telemetryByCode[parameterCode] = latestTelemetry.Value
	}

	// Process each processed parameter
	for parentParamId, parentParam := range listenerFluxor.processedParametersMap {
		sequences, exists := listenerFluxor.processedParameterSequencesMap[parentParamId]
		if !exists || len(sequences) == 0 {
			logger.Warn("No sequences found for processed parameter %d", parentParamId)
			continue
		}

		// Sort sequences by sequence number
		sortedSequences := make([]entity.ProcessedParameterSequence, len(sequences))
		copy(sortedSequences, sequences)
		for i := 0; i < len(sortedSequences); i++ {
			for j := i + 1; j < len(sortedSequences); j++ {
				if sortedSequences[i].Sequence > sortedSequences[j].Sequence {
					sortedSequences[i], sortedSequences[j] = sortedSequences[j], sortedSequences[i]
				}
			}
		}

		// Calculate the processed value
		result := 0.0
		timestamp := time.Now()

		for _, seq := range sortedSequences {

			// Get the value from current telemetry snapshot
			value, exists := telemetryByCode[seq.Parameter.Code]
			if !exists {
				logger.Debug("Parameter %s not found in current telemetry", seq.Parameter.Code)
				continue
			}

			// Apply operation based on Type
			switch seq.Type {
			case "ADDITION":
				result += value
			case "SUBTRACTION":
				result -= value
			case "MULTIPLICATION":
				result *= value
			case "DIVISION":
				if value != 0 {
					result /= value
				} else {
					logger.Warn("Division by zero attempted for processed parameter %d", parentParamId)
					result = 0
				}
			default:
				logger.Warn("Unknown operation type: %s for sequence %d", seq.Type, seq.Id)
			}

			// Update timestamp from latest telemetry
			if tel, ok := listenerFluxor.latestTelemetry[seq.Parameter.Code]; ok {
				timestamp = tel.Timestamp
			}
		}
		if parentParam.IsWatch {
			listenerFluxor.checkAbnormality(parentParam, result)
		}
		// Create telemetry entry for the processed parameter
		processedTelemetry = append(processedTelemetry, &entity.Telemetry{
			ParameterId: parentParam.Id,
			Value:       result,
			Timestamp:   timestamp,
		})

		logger.Debug("Calculated processed parameter %s (ID: %d) = %.2f",
			parentParam.Code, parentParam.Id, result)
	}

	return processedTelemetry
}

func (listenerFluxor *ListenerFluxor) checkAbnormality(parameterEntity *entity.Parameter, value float64) {
	listenerFluxor.alarmMutex.Lock()
	defer listenerFluxor.alarmMutex.Unlock()

	requiredSamples := int(parameterEntity.AbnormalDuration) * 4
	// Init buffer
	if _, isExists := listenerFluxor.abnormalBuffers[parameterEntity.Id]; !isExists {
		listenerFluxor.abnormalBuffers[parameterEntity.Id] = &AbnormalWindow{}
	}

	buffer := listenerFluxor.abnormalBuffers[parameterEntity.Id]
	buffer.Values = append(buffer.Values, value)

	if len(buffer.Values) > requiredSamples {
		buffer.Values = buffer.Values[1:]
	}

	// === CHECK CREATE ALARM ===
	fmt.Println(buffer.Values)
	if len(buffer.Values) == requiredSamples && allAbnormal(buffer.Values, parameterEntity) {
		listenerFluxor.saveAlarm(parameterEntity, value)
		listenerFluxor.recoveryBuffers[parameterEntity.Id] = &AlarmRecoveryWindow{}
		return
	}

	// === CHECK RECOVERY ===
	listenerFluxor.checkRecovery(parameterEntity, value)
}

func allAbnormal(values []float64, parameterEntity *entity.Parameter) bool {
	for _, v := range values {
		if v >= parameterEntity.MinValue && v <= parameterEntity.MaxValue {
			return false
		}
	}
	return true
}
func (listenerFluxor *ListenerFluxor) saveAlarm(parameterEntity *entity.Parameter, value float64) {
	var alarmLogEntity entity.AlarmLog
	err := listenerFluxor.gormDatabase.
		Where("parameter_id = ? AND is_active = ?", parameterEntity.Id, true).
		First(&alarmLogEntity).Error
	eventType := model.AlarmOpen
	isNewAlarm := false // Flag untuk alarm baru

	if err == nil {
		alarmLogEntity.Value = value
		alarmLogEntity.UpdatedAt = time.Now()
	} else {
		isNewAlarm = true // Ini alarm baru
		alarmType := "HIGH"
		if value < parameterEntity.MinValue {
			alarmType = "LOW"
		}

		alarmLogEntity = entity.AlarmLog{
			ParameterId: parameterEntity.Id,
			Value:       value,
			Type:        alarmType,
			IsActive:    true,
			Note:        "",
			Status:      "Open",
		}
	}

	listenerFluxor.gormDatabase.Save(&alarmLogEntity)

	fmt.Println(isNewAlarm)
	// Hanya broadcast jika alarm benar-benar baru
	if isNewAlarm {
		payload, _ := json.Marshal(model.AlarmEvent{
			Id:          alarmLogEntity.Id,
			Type:        eventType,
			ParameterId: parameterEntity.Id,
			Value:       value,
			Status:      alarmLogEntity.Status,
			Timestamp:   time.Now(),
		})
		listenerFluxor.websocketHub.Broadcast(payload)
		logger.Warn("New alarm triggered for parameter %d", parameterEntity.Id)
	}
}
func (listenerFluxor *ListenerFluxor) checkRecovery(parameterEntity *entity.Parameter, value float64) {
	recovery, ok := listenerFluxor.recoveryBuffers[parameterEntity.Id]
	if !ok {
		return
	}

	isNormal := value >= parameterEntity.MinValue && value <= parameterEntity.MaxValue

	if isNormal {
		recovery.NormalCount++
	} else {
		recovery.NormalCount = 0
	}

	if recovery.NormalCount >= 20 {
		listenerFluxor.resolveAlarm(parameterEntity.Id)
		delete(listenerFluxor.recoveryBuffers, parameterEntity.Id)
		delete(listenerFluxor.abnormalBuffers, parameterEntity.Id)
	}
}

func (listenerFluxor *ListenerFluxor) resolveAlarm(parameterId uint64) {
	listenerFluxor.gormDatabase.Model(&entity.AlarmLog{}).
		Where("parameter_id = ? AND is_active = true", parameterId).
		Updates(map[string]interface{}{
			"is_active":   false,
			"status":      "Resolved",
			"resolved_at": time.Now(),
		})

	payload, _ := json.Marshal(model.AlarmEvent{
		Type:        model.AlarmResolved,
		ParameterId: parameterId,
		Status:      "Resolved",
		Timestamp:   time.Now(),
	})
	listenerFluxor.websocketHub.Broadcast(payload)

	logger.Infof("Alarm resolved for parameter %d", parameterId)
}

func startWebSocketServer(hub *websocketConfig.Hub) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketConfig.ServeWS(hub, w, r)
	})

	logger.Info("WebSocket server running on :8181")
	if err := http.ListenAndServe(":8181", nil); err != nil {
		logger.WithError(err).Fatal("Error starting WebSocket server")
	}
}
