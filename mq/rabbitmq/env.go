package mq_serivces

import (
	"fmt"
	"os"
)


func GetEnvWithDefault(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); !ok {
		return defaultVal
	} else {
		return val
	}
}

func GetRabbitEnv() string {
	m := make(map[string]interface{})
	m["host"] = GetEnvWithDefault("RABBIT_MQ_HOST", "192.168.3.101")
	m["port"] = GetEnvWithDefault("RABBIT_MQ_PORT", "5672" )
	m["user"] = GetEnvWithDefault("RABBIT_MQ_USER", "root")
	m["password"] = GetEnvWithDefault("RABBIT_MQ_PASSWD", "root")
	m["vhost"] = GetEnvWithDefault("RABBIT_MQ_VHOST", "/")
	return fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		m["user"],
		m["password"],
		m["host"],
		m["port"],
		m["vhost"])
}





