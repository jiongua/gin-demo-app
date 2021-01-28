package interal

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

func GetDbEnv() map[string]interface{} {
	m := make(map[string]interface{})
	m["host"] = GetEnvWithDefault("DB_HOST", "192.168.3.101")
	m["port"] = GetEnvWithDefault("DB_PORT", "5432" )
	m["user"] = GetEnvWithDefault("DB_USER", "root")
	m["password"] = GetEnvWithDefault("DB_PASSWD", "root")
	m["dbname"] = GetEnvWithDefault("DB_NAME", "zhihudb")

	return m
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

func GetReporterURL() string {
	return GetEnvWithDefault("REPORTER_URL", "http://localhost:8082/reporter/user_action")
}





