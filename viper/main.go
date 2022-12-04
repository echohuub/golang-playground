package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port        int    `mapstructure:"port"`
	Version     string `mapstructure:"version"`
	MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	DbName string `mapstructure:"dbname"`
}

func main() {
	// 设置默认值
	viper.SetDefault("fileDir", "./")
	// 读取配置文件
	viper.SetConfigName("config") // 配置文件名称（无扩展名）
	viper.SetConfigType("yaml")   // 如果配置文件名称里没有扩展名，则需要配置此项
	// viper.SetConfigFile("config.yaml") // 直接配置完整的文件
	viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}

	fmt.Printf("config: %#v\n", viper.AllSettings())
	fmt.Printf("port: %v\n", viper.GetInt("port"))
	fmt.Printf("version: %v\n", viper.GetString("version"))
	fmt.Printf("mysql.host: %v\n", viper.GetString("mysql.host"))
	fmt.Printf("mysql.host: %v\n", viper.Sub("mysql").GetString("host"))

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("viper.Unmarshal falied, err:%v\n", err)
		return
	}
	fmt.Printf("c:%#v\n", c)
}
