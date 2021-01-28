package kafka

import (
	"os"
	"strings"
)

func GetEnvWithDefault(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); !ok {
		return defaultVal
	} else {
		return val
	}
}


func GetKafkaBrokers() []string {
	brokers := GetEnvWithDefault("KAFKA_BROKERS", "192.168.3.101:32774;192.168.3.101:32775")
	return strings.Split(brokers, ",")
}





