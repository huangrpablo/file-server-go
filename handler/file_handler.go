package handler

import (
	"context"
	"errors"
	"github.com/file-server-go/service"
	"github.com/file-server-go/storage"
	"github.com/file-server-go/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"os"
)

// Register routers and their handleFuncs
func (h *FileHandler) Register(r *gin.RouterGroup) {
	r.POST("/files", h.Upload)
	r.GET("/file/:filename", h.Download)
}

type FileHandler struct {
	upload   *service.UploadService
	download *service.DownloadService
}

func New(store storage.FileStore, crypto types.Crypto) *FileHandler {
	upload := service.NewUploadService(store, crypto)
	download := service.NewDownloadService(store, crypto)

	return &FileHandler{
		upload:   upload,
		download: download,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	ctx := c.Request.Context()

	// bind params & validate
	var params UploadParams

	if err := c.ShouldBind(&params); err != nil {
		slog.ErrorContext(ctx, "Get file", "err", err)
		BadRequest(c, "No file provided")
		return
	}

	filename := params.File.Filename

	// to avoid conflicts from concurrent requests
	filepath := "cached/" + filename + uuid.New().String()

	if err := c.SaveUploadedFile(params.File, filepath); err != nil {
		slog.ErrorContext(ctx, "Save file", "err", err)
		InternalServerError(c, "Fail to open the uploaded file")
		return
	}

	defer func() { go h.removeFile(ctx, filepath) }()

	// execute the service
	if err := h.upload.Execute(ctx, filename, filepath); err != nil {
		slog.ErrorContext(ctx, "Upload file", "err", err)
		InternalServerError(c, "Fail to store the uploaded file")
		return
	}

	// return the response
	OK(c, "File uploaded successfully")
}

func (h *FileHandler) Download(c *gin.Context) {
	ctx := c.Request.Context()

	// bind params & validate
	var params DownloadParams

	if err := c.ShouldBindUri(&params); err != nil {
		slog.ErrorContext(ctx, "Get filename", "err", err)
		BadRequest(c, "No filename provided")
		return
	}

	// execute the service
	var (
		filename = params.FileName
		filepath string
		err      error
	)

	if filepath, err = h.download.Execute(ctx, filename); err != nil {
		if errors.Is(err, types.ErrFileNotFound) {
			slog.ErrorContext(ctx, "Retrieve file", "err", err)
			NotFound(c, "Fail to find the requested file")
			return
		}

		slog.ErrorContext(ctx, "Retrieve file", "err", err)
		InternalServerError(c, "Fail to download the requested file")
		return
	}

	defer func() { go h.removeFile(ctx, filepath) }()

	// return the response
	c.File(filepath)
}

func (h *FileHandler) removeFile(ctx context.Context, filepath string) {
	err := os.RemoveAll(filepath)

	if err != nil {
		slog.WarnContext(ctx, "Fail to delete the file from the server", "file", filepath, "err", err)
	}
}
