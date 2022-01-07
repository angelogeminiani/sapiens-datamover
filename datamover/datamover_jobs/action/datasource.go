package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DataMoverDatasource struct {
	name     string
	settings *datamover_commons.DataMoverDatasourceSettings

	_db *gorm.DB
}

func NewDataMoverDatasource(name string, settings *datamover_commons.DataMoverDatasourceSettings) *DataMoverDatasource {
	instance := new(DataMoverDatasource)
	instance.name = name
	instance.settings = settings

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverDatasource) GetData() (interface{}, error) {
	command := instance.settings.Command
	script := instance.settings.Script
	db, err := instance.connection()
	if nil != err {
		return nil, err
	}

	if len(command) > 0 {
		tx := db.Raw(command)
		if nil != tx.Error && !IsRecordNotFoundError(tx.Error) {
			// query error
			err = gg.Errors.Prefix(err, fmt.Sprintf("Error running command '%s': ", command))
			return nil, err
		}
		if tx.RowsAffected > 0 {

		} else {
			// empty data

		}
	}

	if len(script) > 0 {

	}
	return nil, nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverDatasource) connection() (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	if nil == instance._db {
		driver := instance.settings.Connection.Driver
		dsn := instance.settings.Connection.Dsn
		switch driver {
		case "sqlite":
			db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		case "mysql":
			// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		case "postgres":
			// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		case "sqlserver":
			// "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
			db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		}
		if nil == err {
			instance._db = db
		}
	}

	return instance._db, err
}
