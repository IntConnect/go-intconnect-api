package model

type MqttConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	PrefixClient string
	Interval     int
}

type MqttTopicListenerResponse struct {
	SubscribeMultiple SubscribeMultiple
	TopicInstrument   TopicInstrument
}

type TopicInstrument map[string]uint64
type SubscribeMultiple map[string]byte
