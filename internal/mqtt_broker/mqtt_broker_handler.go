package mqtt_broker

import (
	"encoding/json"
	"fmt"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type Handler struct {
	mqttBrokerService Service
	viperConfig       *viper.Viper
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(httpRequest *http.Request) bool { return true },
}

func NewHandler(mqttBrokerService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		mqttBrokerService: mqttBrokerService,
		viperConfig:       viperConfig,
	}
}

func (mqttBrokerHandler *Handler) FindAllMqttBroker(ginContext *gin.Context) {
	mqttBrokerResponses := mqttBrokerHandler.mqttBrokerService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("MqttBroker has been fetched", mqttBrokerResponses))
}
func (mqttBrokerHandler *Handler) GatewayMqttBroker(ginContext *gin.Context) {
	ws, err := upgrader.Upgrade(ginContext.Writer, ginContext.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	var (
		mqttClient mqtt.Client
		mqttLock   sync.Mutex // lock untuk mqttClient
		wsLock     sync.Mutex // ← TAMBAH: lock terpisah untuk ws.WriteMessage
	)

	// Helper: tulis ke WebSocket dengan aman dari goroutine manapun
	writeWS := func(v any) {
		data, err := json.Marshal(v)
		if err != nil {
			return
		}
		wsLock.Lock()
		defer wsLock.Unlock()
		ws.WriteMessage(websocket.TextMessage, data)
	}

	defer func() {
		mqttLock.Lock()
		if mqttClient != nil && mqttClient.IsConnected() {
			mqttClient.Disconnect(250)
		}
		mqttLock.Unlock()
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			return
		}

		var msg model.WebsocketMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(msg)

		switch msg.Type {

		case "connect":
			if msg.BrokerInfo == nil {
				writeWS(map[string]any{"type": "error", "error": "brokerInfo is required"})
				continue
			}

			b := msg.BrokerInfo

			mqttLock.Lock()
			if mqttClient != nil && mqttClient.IsConnected() {
				mqttClient.Disconnect(250)
			}
			mqttLock.Unlock()

			opts := mqtt.NewClientOptions()

			// ✅ FIX 1: gunakan %d bukan %s untuk integer port
			brokerURL := fmt.Sprintf("tcp://%s:%s", b.HostName, b.MqttPort)
			log.Println("[gateway] connecting to broker:", brokerURL)
			opts.AddBroker(brokerURL)

			if b.Username != "" {
				opts.SetUsername(b.Username)
				opts.SetPassword(b.Password)
			}

			opts.SetClientID(fmt.Sprintf("gw-%d", time.Now().UnixNano()))
			fmt.Println(1)
			// ✅ FIX 2: gunakan wsLock saat menulis dari goroutine MQTT
			opts.SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
				log.Println("[gateway] message received:", m.Topic(), string(m.Payload()))
				writeWS(model.WebsocketMessage{
					Type:    "message",
					Topic:   m.Topic(),
					Payload: string(m.Payload()),
				})
			})
			fmt.Println(2)
			client := mqtt.NewClient(opts)
			token := client.Connect()
			token.Wait()

			if token.Error() != nil {
				log.Println("[gateway] connect error:", token.Error())
				writeWS(map[string]any{"type": "error", "error": token.Error().Error()})
				continue
			}

			mqttLock.Lock()
			mqttClient = client
			mqttLock.Unlock()

			log.Println("[gateway] broker connected, notifying browser")
			writeWS(map[string]any{"type": "connected"})

		case "subscribe":
			log.Println("[gateway] subscribe:", msg.Topic)

			mqttLock.Lock()
			c := mqttClient
			mqttLock.Unlock()

			if c == nil || !c.IsConnected() {
				writeWS(map[string]any{"type": "error", "error": "not connected to broker"})
				continue
			}

			token := c.Subscribe(msg.Topic, 0, nil)
			token.Wait()
			if token.Error() != nil {
				log.Println("[gateway] subscribe error:", token.Error())
				writeWS(map[string]any{"type": "error", "error": token.Error().Error()})
			}

		case "publish":
			mqttLock.Lock()
			c := mqttClient
			mqttLock.Unlock()

			if c == nil || !c.IsConnected() {
				continue
			}

			c.Publish(msg.Topic, 0, false, msg.Payload)
		}
	}
}
func (mqttBrokerHandler *Handler) FindAllMqttBrokerPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttBrokerHandler *Handler) CreateMqttBroker(ginContext *gin.Context) {
	var createMqttBrokerModel model.CreateMqttBrokerRequest
	err := ginContext.ShouldBindBodyWithJSON(&createMqttBrokerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.Create(ginContext, &createMqttBrokerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttBrokerHandler *Handler) UpdateMqttBroker(ginContext *gin.Context) {
	var updateMqttBrokerModel model.UpdateMqttBrokerRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateMqttBrokerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	mqttBrokerId := ginContext.Param("id")
	parsedMqttBrokerId, err := strconv.ParseUint(mqttBrokerId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateMqttBrokerModel.Id = parsedMqttBrokerId
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.Update(ginContext, &updateMqttBrokerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttBrokerHandler *Handler) DeleteMqttBroker(ginContext *gin.Context) {
	var deleteMqttBrokerModel model.DeleteResourceGeneralRequest
	mqttBrokerId := ginContext.Param("id")
	parsedMqttBrokerId, err := strconv.ParseUint(mqttBrokerId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))

	err = ginContext.ShouldBindBodyWithJSON(&deleteMqttBrokerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteMqttBrokerModel.Id = parsedMqttBrokerId
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.Delete(ginContext, &deleteMqttBrokerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
