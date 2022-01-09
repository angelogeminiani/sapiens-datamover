package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/schema"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DataMoverDatasource struct {
	root        string
	settings    *datamover_commons.DataMoverDatasourceSettings
	fnVarEngine *gg_fnvars.FnVarsEngine

	_db    *gorm.DB
	schema *schema.DataMoverDatasourceSchema
}

func NewDataMoverDatasource(root string, fnVarEngine *gg_fnvars.FnVarsEngine, datasource *datamover_commons.DataMoverDatasourceSettings) *DataMoverDatasource {
	instance := new(DataMoverDatasource)
	instance.root = root
	instance.fnVarEngine = fnVarEngine
	instance.settings = datasource
	if nil != datasource.Connection {
		instance.schema = datasource.Connection.Schema
	}

	_ = instance.init()

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverDatasource) GetData(context []map[string]interface{}, command string) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	db, err := instance.connection()
	if nil != err {
		// query error
		err = gg.Errors.Prefix(err, fmt.Sprintf("Error opening database '%s': ",
			instance.settings.Connection.Dsn))
		return nil, err
	}

	if nil == context {
		// solve fnvars
		fnRes, fnErr := instance.fnVarEngine.SolveText(command)
		if nil == fnErr {
			command = fnRes
		}
		result, err = query(db, command, nil)
		if nil != err {
			return nil, err
		}
	} else {
		for _, data := range context {
			// solve fnvars
			fnRes, fnErr := instance.fnVarEngine.SolveText(command, data)
			if nil == fnErr {
				command = fnRes
			}
			statement := ToSQLStatement(command, data)
			r, e := query(db, statement, data)
			if nil != e {
				return nil, e
			}
			if len(r) == 0 {
				result = append(result, data)
			} else {
				result = append(result, r...)
			}
		}
	}
	return result, nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

// initialize the schema
func (instance *DataMoverDatasource) init() error {
	if nil != instance {
		db, err := instance.connection()
		if nil != err {
			return err
		}
		if nil == instance.schema {
			// AUTO-CREATE SCHEMA
			instance.schema = schema.NewSchema()
			tbls, err := db.Migrator().GetTables()
			if nil != err {
				return err
			}
			for _, tbl := range tbls {
				tx := db.Table(tbl)
				s := &struct{}{}
				cols, e := tx.Migrator().ColumnTypes(&s)
				if nil != e {
					return e
				}
				// add table to schema
				table := schema.NewTable()
				table.Name = tbl
				for _, col := range cols {
					column := &schema.DataMoverDatasourceSchemaColumn{
						Name: col.Name(),
					}
					column.Nullable, _ = col.Nullable()
					column.Type = col.DatabaseTypeName()
					table.Columns = append(table.Columns, column)
				}
				instance.schema.Tables = append(instance.schema.Tables, table)
				// fmt.Println(table)
			}
		} else {
			// AUTO-MIGRATE SCHEMA
			tables := instance.schema.Tables
			for _, table := range tables {
				if !db.Migrator().HasTable(table.Name) {
					tx := db.Raw("CREATE TABLE ?", table.Name)
					if nil != tx.Error {
						return tx.Error
					}
				}
				s := table.Struct()
				e := db.Table(table.Name).Migrator().AutoMigrate(s)
				if nil != e {
					return e
				}
			}
		}
	}
	return nil
}

// retrieve a connection
func (instance *DataMoverDatasource) connection() (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	if nil == instance._db {
		driver := instance.settings.Connection.Driver
		dsn := instance.settings.Connection.Dsn
		switch driver {
		case "sqlite":
			filename := gg.Paths.Concat(instance.root, dsn)
			db, err = gorm.Open(sqlite.Open(filename), &gorm.Config{})
		case "mysql":
			// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		case "postgres":
			// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		case "sqlserver":
			// "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
			db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		default:
			db = nil
			err = gg.Errors.Prefix(datamover_commons.DatabaseNotSupportedError,
				fmt.Sprintf("'%s': ", driver))
		}
		if nil == err {
			instance._db = db
		}
	}

	return instance._db, err
}
