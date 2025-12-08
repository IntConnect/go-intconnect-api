package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-intconnect-api/pkg/logger"
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
	NumInsertionWorkers = 25
	InsertionQueueSize  = 1000
)

func main() {
	contextWithCancel, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	listenerFluxorInstance := NewListenerFluxor()
	listenerFluxorInstance.StartPeriodicChecker()
	listenerFluxorInstance.StartWorkers(contextWithCancel)
	listenerFluxorInstance.StartTopicHandler(contextWithCancel)
	listenerFluxorInstance.StartSnapshotSaver(contextWithCancel)
	listenerFluxorInstance.WaitForShutdown()
}

type ListenerFluxor struct {
	gormDatabase         *gorm.DB
	mqttBrokersMap       map[uint64]entity.MqttBroker
	parametersMap        map[string]entity.Parameter
	insertionChan        chan []*entity.Telemetry
	waitGroup            *sync.WaitGroup
	telemetryRepository  telemetry.Repository
	parameterRepository  parameter.Repository
	mqttBrokerRepository mqttBroker.Repository
	mqttTopicRepository  mqttTopic.Repository
	mqttClient           mqtt.Client
	rwMutex              sync.RWMutex
	latestTelemetry      map[string]*entity.Telemetry
	telemetryMutex       sync.Mutex
}

func NewListenerFluxor() *ListenerFluxor {
	viperConfig := injector.NewViperConfig()
	databaseCredentials := injector.NewDatabaseCredentials(viperConfig)
	gormDatabase := injector.NewDatabaseConnection(databaseCredentials)
	telemetryRepository := telemetry.NewRepository()
	parameterRepository := parameter.NewRepository()
	mqttBrokerRepository := mqttBroker.NewRepository()
	mqttTopicRepository := mqttTopic.NewRepository()
	latestTelemetry := make(map[string]*entity.Telemetry)
	listenerFluxor := &ListenerFluxor{
		gormDatabase:         gormDatabase,
		mqttBrokersMap:       make(map[uint64]entity.MqttBroker),
		parametersMap:        make(map[string]entity.Parameter),
		insertionChan:        make(chan []*entity.Telemetry, InsertionQueueSize),
		waitGroup:            &sync.WaitGroup{},
		telemetryRepository:  telemetryRepository,
		parameterRepository:  parameterRepository,
		mqttBrokerRepository: mqttBrokerRepository,
		mqttTopicRepository:  mqttTopicRepository,
		latestTelemetry:      latestTelemetry,
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
	parameterEntities, err := listenerFluxor.parameterRepository.FindAll(listenerFluxor.gormDatabase)
	if err != nil {
		return err
	}

	listenerFluxor.rwMutex.Lock()
	defer listenerFluxor.rwMutex.Unlock()
	for _, mqttBrokerEntity := range mqttBrokerEntities {
		listenerFluxor.mqttBrokersMap[mqttBrokerEntity.Id] = mqttBrokerEntity
	}
	for _, parameterEntity := range parameterEntities {
		listenerFluxor.parametersMap[parameterEntity.Code] = parameterEntity
	}
	logger.Infof("Loaded %d mqtt brokers and %d parameters", len(listenerFluxor.mqttBrokersMap), len(listenerFluxor.parametersMap))
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
		ticker := time.NewTicker(1 * time.Minute)
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
	parameterEntities, err := listenerFluxor.parameterRepository.FindAll(listenerFluxor.gormDatabase)
	if err != nil {
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return
	}

	newMqttBrokersMap := make(map[uint64]entity.MqttBroker)
	newParametersMap := make(map[string]entity.Parameter)

	for _, mqttBrokerEntity := range mqttBrokerEntities {
		newMqttBrokersMap[mqttBrokerEntity.Id] = mqttBrokerEntity
	}
	for _, parameterEntity := range parameterEntities {
		newParametersMap[parameterEntity.Code] = parameterEntity
	}

	listenerFluxor.rwMutex.RLock()
	changed := isMqttBrokerMapChanged(listenerFluxor.mqttBrokersMap, newMqttBrokersMap)
	listenerFluxor.rwMutex.RUnlock()

	if changed {
		logger.Warn("MQTT broker configuration changed. Restarting MQTT broker...")
		listenerFluxor.rwMutex.Lock()
		listenerFluxor.mqttBrokersMap = newMqttBrokersMap
		listenerFluxor.parametersMap = newParametersMap
		listenerFluxor.rwMutex.Unlock()

		listenerFluxor.RestartMqttBroker()
	} else {
		listenerFluxor.rwMutex.Lock()
		listenerFluxor.parametersMap = newParametersMap
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
	for i := 1; i <= NumInsertionWorkers; i++ {
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
		if !isExists {
			newParameters = append(newParameters, &entity.Parameter{
				Name:        mqttKey,
				Code:        mqttKey,
				Unit:        "",
				MinValue:    0,
				MaxValue:    0,
				Description: "",
				MqttTopicId: detailMqttTopic["mqtt_topic_id"],
				MachineId:   detailMqttTopic["machine_id"],
			})
		}
	}

	if len(newParameters) > 0 {
		err := listenerFluxor.parameterRepository.CreateBatch(listenerFluxor.gormDatabase, newParameters)
		if err != nil {
			logger.WithError(err).Error("Error creating new parameters")
		}
		parameterEntities, err := listenerFluxor.parameterRepository.FindAll(listenerFluxor.gormDatabase)
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

func isMqttBrokerMapChanged(oldMap, newMap map[uint64]entity.MqttBroker) bool {
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
		snapshotTicker := time.NewTicker(1 * time.Minute)
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
	snapshotPayload := make([]*entity.Telemetry, 0, len(listenerFluxor.latestTelemetry))
	for _, t := range listenerFluxor.latestTelemetry {
		snapshotPayload = append(snapshotPayload, t)
	}
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
