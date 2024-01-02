package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Configuration struct {
	Minio struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
	} `yaml:"minio"`

	App struct {
		ChunkSize  uint64 `yaml:"chunkSize"`
		BucketName string `yaml:"bucketName"`
	} `yaml:"app"`
}

func Load() (*Configuration, error) {
	rawConf, err := os.ReadFile("config/config.yaml")

	if err != nil {
		return nil, fmt.Errorf("fail to read config file: %w", err)
	}

	var conf Configuration
	err = yaml.Unmarshal(rawConf, &conf)

	if err != nil {
		return nil, fmt.Errorf("fail to decode config file: %w", err)
	}

	return &conf, nil
}
