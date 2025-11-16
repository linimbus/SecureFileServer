package main

import (
	"encoding/json"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

type LogConfig struct {
	Filename string `json:"filename"` // 日志文件名
	Level    int    `json:"level"`    // 日志级别
	MaxLines int    `json:"maxlines"` // 最大行数
	MaxSize  int    `json:"maxsize"`  // 最大文件大小
	Daily    bool   `json:"daily"`    // 是否每天生成新日志
	MaxDays  int    `json:"maxdays"`  // 日志保留天数
	Color    bool   `json:"color"`    // 是否彩色输出
}

func LogInit(logPath string) error {
	var logConfig = LogConfig{
		Filename: "debug.log",
		Level:    logs.LevelInformational,
		Daily:    true,
		MaxSize:  100 * 1024 * 1024,
		MaxLines: 100 * 1024,
		MaxDays:  30,
		Color:    false,
	}

	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	if logPath != "" {
		err := EnsureDir(logPath)
		if err != nil {
			logs.Error("ensureDir failed: %v", err)
			return err
		}
		logConfig.Filename = filepath.Join(logPath, logConfig.Filename)
		value, err := json.Marshal(&logConfig)
		if err != nil {
			logs.Error("json.Marshal failed: %v", err)
			return err
		}
		err = logs.SetLogger(logs.AdapterFile, string(value))
		if err != nil {
			logs.Error("logs.SetLogger failed: %v", err)
			return err
		}
		logs.Info("Log Path: %s", AppConf.LogPath)
	}

	return nil
}
