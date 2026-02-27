package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	Level int    `yaml:"level"`
	Path  string `yaml:"path"`
}
type MysqlConfig struct {
	Master   string `yaml:"master"`
	Slave    string `yaml:"slave"`
	LogLevel int    `yaml:"log_level"`
	LogPath  string `yaml:"log_path"`
}
type Config struct {
	Env           string      `yaml:"env"`
	Log           LogConfig   `yaml:"log"`
	MysqlSchedule MysqlConfig `yaml:"mysql_schedule"`
	MysqlPaytrade MysqlConfig `yaml:"mysql_paytrade"`
	MysqlPaycore  MysqlConfig `yaml:"mysql_paycore"`
	MysqlNotify   MysqlConfig `yaml:"mysql_notify"`
	MysqlBilling  MysqlConfig `yaml:"mysql_billing"`
	MysqlCommon   MysqlConfig `yaml:"mysql_common"`
	Server        struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Interval struct {
		Load  int `yaml:"load"`
		Retry int `yaml:"retry"`
		Exec  int `yaml:"exec"`
	} `yaml:"interval"`
}

// GetConfig 获取配置
func GetConfig() *Config {

	return conf
}

var conf *Config

func LoadConfig(path string) {
	fmt.Println("loadConfig, path:", path)
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Println("Config file loaded successfully")

	// 初始化conf变量
	conf = &Config{}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		panic(err)
	}

}

func GetBaseUrl() string {
	if conf.Env == "prod" {
		return "http://pay-scheduler.alphicloud.com"
	} else if conf.Env == "test" {
		return "http://pay-scheduler.hektarapp.io"
	} else {
		return "http://localhost:8080"
	}
}
