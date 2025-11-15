package main

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type DownloadController struct {
	beego.Controller
}

func (c *DownloadController) DownloadFile() {
	uuid := c.Ctx.Input.Param(":uuid")

	fileRecord, err := QueryFileRecord(uuid)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{
			"code": 404,
			"msg":  "File not found or expired",
		}
		c.ServeJSON()

		logs.Info("File record not found or expired, UUID=%s", uuid)
		return
	}

	cryptData, err := ReadFile(fileRecord.FilePath)
	if err != nil {
		DeleteFileRecord(uuid)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{
			"code": 404,
			"msg":  "File not found",
		}
		c.ServeJSON()

		logs.Info("File not found on disk, UUID=%s, filepath=%s", uuid, fileRecord.FilePath)
		return
	}

	decryptedData, err := Decrypt(cryptData, uuid)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "Failed to decrypt file",
		}
		c.ServeJSON()

		logs.Error("Failed to decrypt file, UUID=%s, error=%v", uuid, err)
		return
	}

	c.Ctx.Output.Header("Content-Type", fileRecord.ContentType)
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename=\""+fileRecord.FileName+"\"")
	c.Ctx.Output.Header("Content-Length", strconv.FormatInt(fileRecord.Size, 10))

	c.Ctx.Output.Body(decryptedData)

	logs.Info("File downloaded successfully: UUID=%s, filename=%s", uuid, fileRecord.FileName)
}
