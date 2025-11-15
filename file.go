package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func init() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cleanup()
		}
	}()
}

func GenerateUint32() uint32 {
	var b [4]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func GenerateFileName() string {
	return fmt.Sprintf("%s_%d", time.Now().Format("20060102150405"), GenerateUint32())
}

func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:4]) + "-" + hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" + hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:])
}

func ReadFile(path string) ([]byte, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		logs.Info("read file failed, %s", err.Error())
	}
	return body, err
}

func WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		logs.Info("write file failed, %s", err.Error())
	}
	return err
}

func EnsureDir(path string) error {
	err := os.MkdirAll(path, 0744)
	if err != nil {
		logs.Info("create directory failed, %s", err.Error())
	}
	return err
}

func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		logs.Info("remove file failed, %s", err.Error())
	}
	return err
}

func cleanup() {
	o := orm.NewOrm()
	var expiredFiles []FileRecord
	_, err := o.QueryTable("file_record").Filter("is_deleted", false).Filter("expired_at__lte", time.Now()).All(&expiredFiles)
	if err != nil {
		logs.Error("Failed to query expired files: %v", err)
		return
	}

	for _, file := range expiredFiles {
		err := RemoveFile(file.FilePath)
		if err != nil {
			logs.Error("Failed to remove expired file from disk: %v", err)
			continue
		}

		file.IsDeleted = true
		_, err = o.Update(&file, "IsDeleted")
		if err != nil {
			logs.Error("Failed to mark file as deleted in database: %v", err)
			continue
		}

		logs.Info("Expired file cleaned up: UUID=%s, filepath=%s", file.UUID, file.FilePath)
	}
}
