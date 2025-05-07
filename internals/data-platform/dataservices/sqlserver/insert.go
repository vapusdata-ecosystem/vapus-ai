package sqlserver

import (
	"context"
	// "database/sql"
	"fmt"
	// "log"
	"reflect"
	"strconv"
	"strings"
)

func (m *SqlServerStore) prepareMysqlInsertStmt(Ctx context.Context, table string, dataset *map[string]any, bulk bool) *string {
	// insert data set into mysql
	fields := ""
	values := ""
	for key, value := range *dataset {
		fields += key + ","
		valType := reflect.TypeOf(value)
		if valType.Kind() == reflect.Ptr {
			valType = valType.Elem()
		}
		switch valType.Kind() {
		case reflect.String:
			values += "'" + value.(string) + "',"
		case reflect.Int:
			values += strconv.Itoa(value.(int)) + ","
		case reflect.Int8:
			values += strconv.Itoa(int(value.(int8))) + ","
		case reflect.Int16:
			values += strconv.Itoa(int(value.(int16))) + ","
		case reflect.Int32:
			values += strconv.Itoa(int(value.(int32))) + ","
		case reflect.Int64:
			values += strconv.Itoa(int(value.(int64))) + ","
		case reflect.Float32:
			values += strconv.FormatFloat(value.(float64), 'f', -1, 64) + ","
		case reflect.Float64:
			values += strconv.FormatFloat(value.(float64), 'f', -1, 64) + ","
		case reflect.Bool:
			values += strconv.FormatBool(value.(bool)) + ","
		default:
			continue // TO:DO: handle other types and default types as well
		}
	}
	fields = strings.TrimSuffix(fields, ",")
	values = strings.TrimSuffix(values, ",")
	stmt := fmt.Sprintf("INSERT INTO `"+table+"` (%v) VALUES (%v)", fields, values)
	return &stmt
}
