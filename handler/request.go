package handler

import "mime/multipart"

type UploadParams struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type DownloadParams struct {
	FileName string `uri:"filename" binding:"required"`
}
