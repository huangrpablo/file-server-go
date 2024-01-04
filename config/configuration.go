package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Configuration struct {
	// Minio credentials
	Minio struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
	} `yaml:"minio"`

	App struct {
		// the size of a chunk if multipart upload is enabled
		ChunkSize uint64 `yaml:"chunkSize"`
		// whether a file will be uploaded to minio in parts
		UploadInChunk bool `yaml:"uploadInChunk"`
		// the bucket where files are stored
		BucketName string `yaml:"bucketName"`
	} `yaml:"app"`
}

func Load(path string) (*Configuration, error) {
	//rawConf, err := os.ReadFile("config/config.yaml")

	rawConf, err := os.ReadFile(path)

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
