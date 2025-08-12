package configs

import (
	"fmt"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/helper"
	"log"
	"strconv"
)

type Mqtt struct {
	Host         string
	Port         string
	Username     string
	Password     string
	PrefixClient string
}

func (mqttConfig *Mqtt) createClientOptions() *pahomqtt.ClientOptions {
	opts := pahomqtt.NewClientOptions()
	port, _ := strconv.Atoi(mqttConfig.Port)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttConfig.Host, port))
	opts.SetUsername(mqttConfig.Username)
	opts.SetPassword(mqttConfig.Password)
	opts.SetClientID(fmt.Sprintf("%s-%v", mqttConfig.PrefixClient, helper.RandomStringGenerator(5)))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	return opts
}

var messagePubHandler pahomqtt.MessageHandler = func(client pahomqtt.Client, msg pahomqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler pahomqtt.OnConnectHandler = func(client pahomqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler pahomqtt.ConnectionLostHandler = func(client pahomqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func (mqttConfig *Mqtt) Connect() pahomqtt.Client {
	opts := mqttConfig.createClientOptions()
	client := pahomqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return client
}

func InitMqtt(mqttConfig *model.MqttConfig) *Mqtt {
	return &Mqtt{
		Host:         mqttConfig.Host,
		Port:         mqttConfig.Port,
		Username:     mqttConfig.Username,
		Password:     mqttConfig.Password,
		PrefixClient: mqttConfig.PrefixClient,
	}
}

func (mqttConfig *Mqtt) PublishToParameter(parameter string, payload string) error {
	client := mqttConfig.Connect()
	defer client.Disconnect(250)

	token := client.Publish(parameter, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("gagal publish ke topic %s: %w", parameter, token.Error())
	}
	return nil
}
