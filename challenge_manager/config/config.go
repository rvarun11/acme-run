package config

import "os"

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// The following can be used to set global config
// var Config *GlobalConfig

// type GlobalConfig struct {
// 	MODE string
// }

// func init() {
// 	Config = newConfig()
// }

// func newConfig() *GlobalConfig {
// 	return &GlobalConfig{
// 		MODE: GetEnv("APP_MODE", "development"),
// 	}
// }
