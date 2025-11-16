package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type FileRecord struct {
	ID          int    `orm:"auto"`
	UUID        string `orm:"size(36);unique"`
	FileName    string `orm:"size(255)"`
	FilePath    string `orm:"size(1024)"`
	Size        int64
	ContentType string    `orm:"size(100)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	ExpiredAt   time.Time
	IsDeleted   bool `orm:"default(false)"`
}

func init() {
	orm.RegisterModel(new(FileRecord))
}

func InitDB() error {
	err := orm.RegisterDriver("sqlite", orm.DRSqlite)
	if err != nil {
		logs.Error("init db failed, %s", err.Error())
		return err
	}
	err = orm.RegisterDataBase("default", "sqlite3", AppConf.DataBasePath)
	if err != nil {
		logs.Error("init db failed, %s", err.Error())
		return err
	}
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		logs.Error("init db failed, %s", err.Error())
		return err
	}
	return nil
}

func CreateFileRecord(file *FileRecord) error {
	o := orm.NewOrm()
	_, err := o.Insert(file)
	if err != nil {
		logs.Error("create file record failed, %s", err.Error())
	}
	return err
}

func QueryFileRecord(uuid string) (*FileRecord, error) {
	o := orm.NewOrm()
	var file FileRecord
	err := o.QueryTable("file_record").Filter("uuid", uuid).Filter("is_deleted", false).Filter("expired_at__gt", time.Now()).One(&file)
	if err != nil {
		logs.Error("get file record %s failed, %s", uuid, err.Error())
		return nil, err
	}
	return &file, nil
}

func DeleteFileRecord(uuid string) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("file_record").Filter("uuid", uuid).Update(orm.Params{
		"is_deleted": true,
	})

	if err != nil {
		logs.Error("delete file record failed, %s", err.Error())
	}
	return err
}
