package handler

import (
	"errors"
	"github.com/file-server-go/service"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

type FileHandler struct {
	uploadSrv   *service.UploadService
	downloadSrv *service.DownloadService
}

func New(store storage.Store, crypto types.Crypto, bucketName string) *FileHandler {
	uploadSrv := service.NewUploadService(store, crypto, bucketName)
	downloadSrv := service.NewDownloadService(store, crypto, bucketName)

	return &FileHandler{
		uploadSrv:   uploadSrv,
		downloadSrv: downloadSrv,
	}
}

func (h *FileHandler) Register(r *gin.Engine) {
	r.POST("/files", h.Upload)
	r.GET("/file/:filename", h.Download)
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		slog.Error("Get file: %+v", err)
		BadRequest(c, "No file provided")
		return
	}

	filepath := "uploaded/" + file.Filename
	err = c.SaveUploadedFile(file, filepath)

	if err != nil {
		slog.Error("Save file: %+v", err)
		InternalServerError(c, "Fail to open the uploaded file")
		return
	}

	defer func() { go h.removeFile(filepath) }()

	err = h.uploadSrv.Upload(file.Filename, filepath)

	if err != nil {
		slog.Error("Upload file: %+v", err)
		InternalServerError(c, "Fail to store the uploaded file")
		return
	}

	OK(c, "File uploaded successfully")
}

func (h *FileHandler) Download(c *gin.Context) {
	filename := c.Param("filename")

	filepath, err := h.downloadSrv.Download(filename)

	if err != nil && errors.Is(err, types.ErrFileNotFound) {
		slog.Error("Retrieve file: %+v", err)
		NotFound(c, "Fail to find the requested file")
		return
	}

	if err != nil {
		slog.Error("Retrieve file: %+v", err)
		InternalServerError(c, "Fail to download the requested file")
		return
	}

	defer func() { go h.removeFile(filepath) }()

	c.File(filepath)
}

func (h *FileHandler) removeFile(filepath string) {
	err := os.RemoveAll(filepath)

	if err != nil {
		slog.Warn("Fail to delete the file[%v] on the server: %+v", filepath, err)
	}
}
