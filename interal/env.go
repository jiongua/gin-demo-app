package interal

import "os"

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

