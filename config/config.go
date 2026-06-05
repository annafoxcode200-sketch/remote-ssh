package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// 定义与 YAML 结构严格对应的结构体
type Config struct {
	Server ServerConfig `yaml:"server"`
	AliECS AliECSConfig `yaml:"ali-ecs"`
}

type ServerConfig struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type AliECSConfig struct {
	Endpoint        string `yaml:"endpoint"`
	InstanceId      string `yaml:"instanceId"`
	AccessKeyID     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
}

// Load 加载配置文件并返回配置对象
func Load(filePath string) *Config {
	// 1. 读取 YAML 文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("解析 YAML 失败: %v", err)
	}

	// 3. 格式化并完整输出一次配置项
	printConfig(&cfg)

	return &cfg
}

// printConfig 将配置序列化为美观的 YAML 格式并打印
func printConfig(cfg *Config) {
	fmt.Println("========== 配置加载完成 ==========")
	// 将结构体重新序列化为 YAML 字节流，这样能保证输出的格式非常整齐且易读
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		log.Printf("序列化配置用于打印时出错: %v", err)
		return
	}
	fmt.Println(string(yamlData))
	fmt.Println("==================================")
}
