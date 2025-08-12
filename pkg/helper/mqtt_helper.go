package helper

import (
	"fmt"
	"strconv"
)

func GetTopicKeys(topics map[string]byte) []string {
	keys := make([]string, 0, len(topics))
	for key := range topics {
		keys = append(keys, key)
	}
	return keys
}

func ExtractValue(rawValue interface{}) (string, error) {
	switch v := rawValue.(type) {
	case float64:
		return fmt.Sprintf("%f", v), nil
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case []interface{}:
		// Jika berupa array/slice dan memiliki satu elemen, ambil elemen pertama
		if len(v) == 1 {
			return ExtractValue(v[0])
		}
		// Jika lebih dari satu, bisa digabung atau dikembalikan error sesuai kebutuhan
		return "", fmt.Errorf("array contains multiple values")

	default:
		return "", fmt.Errorf("unknown data type")
	}
}
