package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/schema"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/scripting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DataMoverDatasource struct {
	root               string
	fnVarEngine        *gg_fnvars.FnVarsEngine
	connectionSettings *datamover_commons.DataMoverConnectionSettings
	scriptContext      string

	_db          *gorm.DB
	schema       *schema.DataMoverDatasourceSchema
	scriptEngine *scripting.ScriptController
}

func NewDataMoverDatasource(root string, fnVarEngine *gg_fnvars.FnVarsEngine,
	connection *datamover_commons.DataMoverConnectionSettings,
	scripts *datamover_commons.DataMoverActionScriptSettings) (*DataMoverDatasource, error) {

	instance := new(DataMoverDatasource)
	instance.root = root
	instance.fnVarEngine = fnVarEngine
	instance.connectionSettings = connection
	if nil != instance.connectionSettings {
		instance.schema = instance.connectionSettings.Schema
	}
	if nil != scripts {
		instance.scriptContext = scripts.Context
	}
	err := instance.init()

	return instance, err
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverDatasource) GetSchema() *schema.DataMoverDatasourceSchema {
	if nil != instance {
		return instance.schema
	}
	return nil
}

func (instance *DataMoverDatasource) GetData(command string, context []map[string]interface{}, variables map[string]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	db, err := instance.connection()
	if nil != err {
		// query error
		err = gg.Errors.Prefix(err, fmt.Sprintf("Error opening database '%s': ",
			instance.connectionSettings.Dsn))
		return nil, err
	}

	if nil == context {
		// solve fnvars
		fnRes, fnErr := instance.fnVarEngine.SolveText(command)
		if nil == fnErr {
			command = fnRes
		}
		result, err = query(db, command, variables)
		if nil != err {
			return nil, err
		}
		// context script
		var scriptVars map[string]interface{}
		result, scriptVars = instance.scriptEngine.RunWithArray("context", instance.scriptContext, result, variables)
		if len(scriptVars) > 0 {
			gg.Maps.Merge(true, variables, scriptVars)
		}
	} else {
		// context script
		var scriptVars map[string]interface{}
		context, scriptVars = instance.scriptEngine.RunWithArray("context", instance.scriptContext, context, variables)
		if len(scriptVars) > 0 {
			gg.Maps.Merge(true, variables, scriptVars)
		}
		for _, data := range context {
			if nil != data {
				ctx := gg.Maps.Merge(false, map[string]interface{}{}, data, variables)
				// solve fnvars
				fnRes, fnErr := instance.fnVarEngine.SolveText(command, ctx)
				if nil == fnErr {
					command = fnRes
				}
				statement := ToSQLStatement(command, data)
				r, e := query(db, statement, ctx)
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
	}
	return result, nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

// initialize the schema, network, etc...
func (instance *DataMoverDatasource) init() error {
	if nil != instance {
		err := instance.initSchema()
		if nil != err {
			return err
		}

		instance.scriptEngine = scripting.NewScriptController(instance.root)
	}
	return nil
}

func (instance *DataMoverDatasource) initSchema() error {
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
	return nil
}

// retrieve a connection
func (instance *DataMoverDatasource) connection() (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	if nil == instance._db {
		driver := instance.connectionSettings.Driver
		dsn := instance.connectionSettings.Dsn
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
