package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

func IsRecordNotFoundError(err error) bool {
	return nil != err && err.Error() == "record not found"
}

// ToSQLStatement convert "INSERT INTO table (...)" into "INSERT INTO table (field1, field2,...) VALUES (value1, value2,...)"
func ToSQLStatement(source string, data map[string]interface{}, fieldsMapping map[string]interface{}) string {
	if strings.Index(source, "(...)") > -1 {
		var fields, values string
		for k, v := range data {
			if len(fields) > 0 {
				fields += ","
				values += ","
			}
			// field
			if n := gg.Maps.Get(fieldsMapping, k); nil != n {
				fields += gg.Convert.ToString(n)
			} else {
				fields += k
			}
			// value
			values += toSQLString(v)
		}
		statement := fmt.Sprintf("(%s) VALUES (%s)", fields, values)
		return strings.ReplaceAll(source, "(...)", statement)
	} else {
		return source
	}
}

// QueryGetParamNames return unique param names
func QueryGetParamNames(query string) []string {
	return getParamNames(query, "@")
}

func LoadJsDatasets(root string) (response map[string][]interface{}) {
	response = make(map[string][]interface{})
	filename := gg.Paths.Concat(root, "datasets.json")
	if ok, _ := gg.Paths.Exists(filename); ok {
		_ = gg.JSON.ReadFromFile(filename, &response)
	}
	return
}

func OverwriteJsDatasets(root string, datasets map[string][]interface{}) {
	filename := gg.Paths.Concat(root, "datasets.json")
	_ = gg.Paths.Mkdir(filename)
	_, _ = gg.IO.WriteTextToFile(gg.JSON.Stringify(datasets), filename)
}

func getParamNames(query, prefix string) []string {
	response := make([]string, 0)
	query = strings.ReplaceAll(query, ";", " ;")
	query += " "
	params := gg.Regex.TextBetweenStrings(query, prefix, " ")
	for _, param := range params {
		// purge name from comma or other invalid delimiters
		param = strings.TrimRight(param, ",.;:\n\r")
		if gg.Arrays.IndexOf(param, response) == -1 {
			response = append(response, param)
		}
	}
	return response
}

func QueryGetNamedArgs(command string, context map[string]interface{}) []interface{} {
	args := make([]interface{}, 0)

	// parse command to get sql parameters
	params := QueryGetParamNames(command)
	if len(params) > 0 {
		for _, param := range params {
			if v, b := context[param]; b {
				args = append(args, sql.Named(param, v))
			}
		}
	}

	return args
}

func toSQLString(v interface{}) string {
	if nil == v {
		return "NULL"
	}
	if dt, ok := v.(time.Time); ok {
		// ISO-8601
		//result := dt.UTC().Format("2006-01-02T15:04:05-0700")
		// yyyy-MM-dd HH:mm:ss
		result := gg.Formatter.FormatDate(dt, "yyyy-MM-dd HH:mm:ss")
		return fmt.Sprintf("'%s'", result)
	}
	val := gg.Reflect.ValueOf(v)
	kind := val.Kind()
	switch kind {
	case reflect.String:
		if s, b := val.Interface().(string); b {
			return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
		} else {
			return fmt.Sprintf("'%v'", v)
		}
	case reflect.Slice, reflect.Array:
		var buf strings.Builder
		a := gg.Convert.ToArray(v)
		for i, av := range a {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(toSQLString(av))
		}
		return buf.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func query(db *gorm.DB, command string, context map[string]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	// preprocess for @@variables
	command = preProcessSpecials(command, context)

	// get args
	args := QueryGetNamedArgs(command, context)

	// do query
	tx := db.Raw(command, args...)
	tx.Scan(&result)
	if nil != tx.Error && !IsRecordNotFoundError(tx.Error) {
		// query error
		return nil, gg.Errors.Prefix(tx.Error, fmt.Sprintf("Error running command '%s': ",
			command))
	}
	return result, nil
}

func preProcessSpecials(command string, context map[string]interface{}) string {
	specials := getParamNames(command, "@@")
	if len(specials) > 0 {
		// replace
		for _, special := range specials {
			if v, ok := context[special]; ok {
				value := toSQLString(v)
				command = strings.ReplaceAll(command, "@@"+special, value)
			}
		}
	}
	return command
}
