package database

import (
    "github.com/go-xorm/xorm"
    "go.uber.org/zap"
)

func NewDBClient(dsn string, mode string, logger *zap.Logger) *MySQLClient {
    client := new(MySQLClient)
    client.logger = logger
    engine, err := xorm.NewEngine("mysql", dsn)
    if err != nil {
        panic(err)
    }
    if mode == "dev" {
        engine.ShowSQL(true)
    }
    client.engine = engine
    return client
}
