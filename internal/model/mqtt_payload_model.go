package model

import "time"

type MqttPayload struct {
	MqttInnerPayload map[string][]interface{} `json:"d"`
	Timestamp        time.Time                `json:"ts"`
}
