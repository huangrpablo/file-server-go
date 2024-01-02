package main

import (
	"github.com/file-server-go/config"
	"github.com/file-server-go/handler"
	"github.com/file-server-go/storage/minio"
	"github.com/file-server-go/types"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func main() {
	conf, err := config.Load()

	if err != nil {
		panic(err)
	}

	slog.Info("init config successfully")

	store, err := minio.New(conf.Minio.Endpoint, conf.Minio.AccessKey, conf.Minio.SecretKey, conf.App.ChunkSize)

	if err != nil {
		panic(err)
	}

	slog.Info("init minio successfully")

	crypto, err := types.NewAES()

	if err != nil {
		panic(err)
	}

	slog.Info("init crypto artifacts successfully")

	fileHandler := handler.New(store, crypto, conf.App.BucketName)

	r := gin.Default()
	fileHandler.Register(r)

	slog.Info("register routing successfully")

	slog.Info("service is up")

	r.Run(":8080")
}
