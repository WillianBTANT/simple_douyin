package config

import (
	"encoding/json"
	"os"
	"reflect"
	"sync"
)

type Config struct {
	DatabaseURI string `json:"database_url"`
	ServerPort  int    `json:"server_port"`
	ServerHost  string `json:"server_host"`

	JWTSecret string `json:"jwt_secret"`
	JWTExp    int64  `json:"jwt_exp"`

	StaticBaseUrl string `json:"static_base_url"`
}

var (
	defaultConfig = Config{
		DatabaseURI:   "root:604486@tcp(127.0.0.1:3306)/simple_douyin?charset=utf8mb4&parseTime=True&loc=Local", // 请修改为自己的数据库url
		ServerPort:    8080,
		ServerHost:    "0.0.0.0",
		JWTSecret:     "Cd18jPdmI6aUWEkmvZomzltcjhetm9MMn64HDEwFBYDpMlks2fZpZ0nGvWp0G9TuMcQeySuGi_P7jwQA3gh7zQ",
		JWTExp:        60 * 60 * 24,                         // 一天
		StaticBaseUrl: "http://192.168.115.89:8080/static/", // 请修改为静态文件的url
	}
	config     Config
	configLock sync.RWMutex
	filePath   = "./config/config.json"
)

func init() {
	_ = ReadConfigFromFile()
}

// ReadConfigFromFile 从文件中读取配置、配置文件不存在就使用默认的配置
func ReadConfigFromFile() error {
	defer func() {
		_ = SaveConfigToFile()
	}()
	file, err := os.Open(filePath)
	if err != nil {
		config = defaultConfig
		return nil
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := json.NewDecoder(file)
	configLock.Lock()
	defer configLock.Unlock()

	if err = decoder.Decode(&config); err != nil {
		return err
	}

	// 使用反射遍历默认配置，并将空字段替换为默认值
	defaultValue := reflect.ValueOf(defaultConfig)
	configValue := reflect.ValueOf(&config).Elem()

	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		if field.Interface() == reflect.Zero(field.Type()).Interface() {
			defaultField := defaultValue.Field(i)
			field.Set(defaultField)
		}
	}
	return nil
}

func GetConfig() Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func SaveConfigToFile() error {
	configLock.RLock()
	defer configLock.RUnlock()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
