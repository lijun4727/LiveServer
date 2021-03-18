package db_connect

import (
	"LiveServer/common"
	"errors"
	"log"

	"github.com/go-xorm/xorm"
)

var MysqlEngine *xorm.Engine

var (
	MysqlConnectError           = errors.New("mysql connect error")
	MysqlCreateAccountDataError = errors.New("mysql create table account_data error")
)

func InitConnect() error {
	var err error
	MysqlEngine, err = xorm.NewEngine("mysql", common.MysqlConnCmd)
	if err != nil {
		log.Fatalf("failed to connect mysql: %v, connect cmd: %s", err, common.MysqlConnCmd)
		return MysqlConnectError
	}

	return nil
}

func Clean() {
	if MysqlEngine != nil {
		MysqlEngine.Close()
	}
}
