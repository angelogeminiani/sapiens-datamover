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
func ToSQLStatement(source string, m map[string]interface{}) string {
	if strings.Index(source, "(...)") > -1 {
		var fields, values string
		for k, v := range m {
			if len(fields) > 0 {
				fields += ","
				values += ","
			}
			fields += k
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
	response := make([]string, 0)
	query = strings.ReplaceAll(query, ";", " ;")
	query += " "
	params := gg.Regex.TextBetweenStrings(query, "@", " ")
	for _, param := range params {
		if gg.Arrays.IndexOf(param, response) == -1 {
			response = append(response, param)
		}
	}
	return response
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
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v.(string), "'", "''"))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func query(db *gorm.DB, command string, context map[string]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	args := make([]interface{}, 0)

	// parse command to get sql parameters
	params := QueryGetParamNames(command)
	if len(params) > 0 {
		for _, param := range params {
			if v, b := context[param]; b {
				args = append(args, sql.Named(param, v))
				// command = strings.ReplaceAll(command, fmt.Sprintf("@%s", param), "?")
			}
		}
	}

	tx := db.Raw(command, args...)
	tx.Scan(&result)
	if nil != tx.Error && !IsRecordNotFoundError(tx.Error) {
		// query error
		return nil, gg.Errors.Prefix(tx.Error, fmt.Sprintf("Error running command '%s': ",
			command))
	}
	return result, nil
}
