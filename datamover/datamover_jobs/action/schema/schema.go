package schema

import (
	"bitbucket.org/digi-sense/gg-core"
	"fmt"
)

type DataMoverDatasourceSchema struct {
	Tables []*DataMoverDatasourceSchemaTable `json:"tables"`
}
type DataMoverDatasourceSchemaTable struct {
	Name    string                             `json:"name"`
	Columns []*DataMoverDatasourceSchemaColumn `json:"columns"`
}

func (instance *DataMoverDatasourceSchemaTable) String() string {
	return gg.JSON.Stringify(instance)
}

func (instance *DataMoverDatasourceSchemaTable) Struct() interface{} {
	builder := gg.Structs.New()
	for _, column := range instance.Columns {
		tag := column.Tag
		if len(tag) > 0 {
			tag = fmt.Sprintf(`gorm:"%s"`, tag)
		}
		name := column.Name
		t := column.Type
		builder.AddFieldByStringType(name, tag, t)
	}
	return builder.Interface()
}

type DataMoverDatasourceSchemaColumn struct {
	Name     string `json:"name"`
	Nullable bool   `json:"nullable"`
	Type     string `json:"type"`
	Tag      string `json:"tag"`
}

func NewSchema() *DataMoverDatasourceSchema {
	instance := new(DataMoverDatasourceSchema)
	instance.Tables = make([]*DataMoverDatasourceSchemaTable, 0)
	return instance
}

func NewTable() *DataMoverDatasourceSchemaTable {
	instance := new(DataMoverDatasourceSchemaTable)
	instance.Columns = make([]*DataMoverDatasourceSchemaColumn, 0)
	return instance
}
