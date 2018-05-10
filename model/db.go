package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	_ "github.com/go-sql-driver/mysql"

	. "github.com/daoone/hammer/util"
)

var DbEngine *xorm.Engine

func init() {
	GetEngine()
}

func NewEngine() (*xorm.Engine) {
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		GetConfig("database", "username"),
		GetConfig("database", "password"),
		GetConfig("database", "address"),
		GetConfig("database", "port"),
		GetConfig("database", "dbname"),
	)
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		DoLog(err.Error(), "model_db")
		panic(err.Error())
	}

	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, GetConfig("database", "tablePrefix"),)
	engine.SetTableMapper(tbMapper)
	engine.ShowSQL(true)

	return engine
}

func GetEngine() (*xorm.Engine) {
	if DbEngine == nil || DbEngine.Ping() != nil {
		DbEngine = NewEngine()
	}
	return DbEngine
}