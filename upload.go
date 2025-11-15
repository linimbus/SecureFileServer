package main

import (
	"io"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) UploadFile() {
	file, header, err := c.GetFile("file")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "Failed to get uploaded file",
		}
		c.ServeJSON()

		logs.Warn("Failed to get uploaded file: %v", err)
		return
	}
	defer file.Close()

	if header.Size > AppConf.MaxFileSize {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "File too large",
		}
		c.ServeJSON()

		logs.Warn("Uploaded file too large: %d bytes", header.Size)
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "Failed to read file data",
		}
		c.ServeJSON()

		logs.Error("Failed to read file data: %v", err)
		return
	}

	if int64(len(fileData)) != header.Size {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "File size mismatch",
		}
		c.ServeJSON()

		logs.Warn("File size mismatch: expected %d, got %d", header.Size, len(fileData))
		return
	}
	uuid := GenerateUUID()

	fileData, err = Encrypt(fileData, uuid)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "Failed to encrypt data",
		}
		c.ServeJSON()

		logs.Error("Failed to encrypt data: %v", err)
		return
	}

	filePath := filepath.Join(AppConf.UploadPath, GenerateFileName())

	err = WriteFile(filePath, fileData)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "Failed to save file",
		}
		c.ServeJSON()

		logs.Error("Failed to save file: %v", err)
		return
	}

	expiredAt := time.Now().Add(time.Duration(AppConf.FileExpiry) * time.Hour)

	fileRecord := &FileRecord{
		UUID:        uuid,
		FileName:    header.Filename,
		FilePath:    filePath,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		ExpiredAt:   expiredAt,
	}

	err = CreateFileRecord(fileRecord)
	if err != nil {
		RemoveFile(filePath)
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "Failed to save file record",
		}
		c.ServeJSON()

		logs.Error("Failed to save file record: %v", err)
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "File uploaded successfully",
		"data": map[string]interface{}{
			"uuid":         uuid,
			"filename":     header.Filename,
			"size":         header.Size,
			"download_url": "/download/" + uuid,
		},
	}
	c.ServeJSON()

	logs.Info("File uploaded successfully: UUID=%s, filename=%s, size=%d bytes", uuid, header.Filename, header.Size)
}
