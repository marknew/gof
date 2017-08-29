package util

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

type MYSQL_DB struct {
	dbusername string
	dbpassowrd string
	dbhostsip  string
	dbname     string
	DB         *sqlx.DB
}

func NewDB(dbhostsip, dbname, dbusername, dbpassowrd string) *MYSQL_DB {
	return &MYSQL_DB{dbusername, dbpassowrd, dbhostsip, dbname, nil}
}

func (f *MYSQL_DB) Mysql_open() error {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&multiStatements=true",
		f.dbusername, f.dbpassowrd, f.dbhostsip, f.dbname))

	CheckErr(err)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Second * 600)
	db.SetMaxOpenConns(500)
	err = db.Ping()
	if err != nil {
		log.V(2).Infoln("连接mysql失败", f.dbhostsip, f.dbname)
		return err
	}
	log.V(2).Infoln("链接成功", f.dbhostsip, f.dbname)

	f.DB = db
	return err
}

func (f *MYSQL_DB) Mysql_close() { //关闭
	defer f.DB.Close()
}
