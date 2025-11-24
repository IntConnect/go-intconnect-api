package model

import "time"
import "strings"

type MqttPayload struct {
	MqttInnerPayload map[string][]interface{} `json:"d"`
	Timestamp        CustomTime                `json:"ts"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`) // hapus quotes
	if s == "" {
		ct.Time = time.Time{}
		return nil
	}

	// Coba parse microseconds tanpa timezone
	t, err := time.Parse("2006-01-02T15:04:05.999999", s)
	if err != nil {
		// fallback: coba RFC3339Nano
		t, err = time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return err
		}
	}

	ct.Time = t
	return nil
}
