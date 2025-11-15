package main

import (
	"flag"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/mattn/go-sqlite3"
)

var VERSION = "1.0.0"

type Config struct {
	Help        bool
	HTTPPort    int
	HTTPHost    string
	UploadPath  string
	MaxFileSize int64
	FileExpiry  int64
	DBPath      string
}

var AppConf Config

func init() {
	os.Args[0] = "SecureFileServer v" + VERSION

	flag.IntVar(&AppConf.HTTPPort, "port", 8080, "HTTP server port")
	flag.StringVar(&AppConf.HTTPHost, "host", "0.0.0.0", "HTTP server host (IP address to bind)")
	flag.StringVar(&AppConf.UploadPath, "upload-path", "./uploads", "Upload directory path")
	flag.Int64Var(&AppConf.MaxFileSize, "max-file-size", 104857600, "Maximum file size in bytes (default 100MB)")
	flag.Int64Var(&AppConf.FileExpiry, "file-expiry", 24, "File expiry duration in hours")
	flag.StringVar(&AppConf.DBPath, "db-path", "./database.db", "Database file path")
	flag.BoolVar(&AppConf.Help, "help", false, "Show help message")
}

func main() {
	flag.Parse()

	if AppConf.Help {
		flag.Usage()
		return
	}

	if AppConf.MaxFileSize <= 0 {
		logs.Error("Max file size must be greater than 0")
		os.Exit(1)
	}

	if AppConf.FileExpiry <= 0 {
		logs.Error("File expiry must be greater than 0")
		os.Exit(1)
	}

	err := InitDB()
	if err != nil {
		logs.Error("Failed to initialize database: %s", err.Error())
		os.Exit(1)
	}

	err = EnsureDir(AppConf.UploadPath)
	if err != nil {
		logs.Error("Failed to create upload directory: %s", err.Error())
		os.Exit(1)
	}

	logs.Info("Starting server on %s:%d", AppConf.HTTPHost, AppConf.HTTPPort)
	logs.Info("Max file size: %d bytes", AppConf.MaxFileSize)
	logs.Info("File expiry: %v hours", AppConf.FileExpiry)
	logs.Info("Upload path: %s", AppConf.UploadPath)
	logs.Info("Load database: %s", AppConf.DBPath)

	beego.BConfig.Listen.HTTPAddr = AppConf.HTTPHost
	beego.BConfig.Listen.HTTPPort = AppConf.HTTPPort
	beego.BConfig.Listen.EnableHTTP = true

	beego.Router("/upload", &UploadController{}, "post:UploadFile")
	beego.Router("/download/:uuid", &DownloadController{}, "get:DownloadFile")

	beego.Run()
}
