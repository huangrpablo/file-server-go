package main

import (
	"github.com/file-server-go/config"
	"github.com/file-server-go/handler"
	"github.com/file-server-go/storage/minio"
	"github.com/file-server-go/types"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func initConfig() *config.Configuration {
	conf, err := config.Load("config/config.yaml")

	if err != nil {
		panic(err)
	}

	return conf
}

func initMinio(conf *config.Configuration) {
	err := minio.Init(conf.Minio.Endpoint, conf.Minio.AccessKey, conf.Minio.SecretKey)

	if err != nil {
		panic(err)
	}
}

func initCrypto() types.Crypto {
	crypto, err := types.NewAES()

	if err != nil {
		panic(err)
	}

	return crypto
}

func initHandler(conf *config.Configuration, crypto types.Crypto) *handler.FileHandler {
	store := minio.NewFileStore(conf.App.ChunkSize, conf.App.UploadInChunk, conf.App.BucketName)

	return handler.New(store, crypto)
}

func main() {
	conf := initConfig()

	slog.Info("init config successfully")

	initMinio(conf)

	slog.Info("init minio successfully")

	crypto := initCrypto()

	slog.Info("init crypto artifacts successfully")

	fileHandler := initHandler(conf, crypto)

	slog.Info("init handler successfully")

	r := gin.Default()
	v1 := r.Group("/v1")

	fileHandler.Register(v1)

	slog.Info("register routing successfully")

	slog.Info("the service is up")

	r.Run(":8080")
}
