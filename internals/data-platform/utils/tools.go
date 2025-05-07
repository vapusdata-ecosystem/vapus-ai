package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

func BuildDataTableName(tableName string) string {
	ot := strings.ToLower(tableName)
	rep := strings.NewReplacer(" ", "_", "+", "_", "-", "_", ".", "_", "/", "_")
	return rep.Replace(ot)
}

func validateType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != expected {
		return nil, ErrInvalidDestinationType
	}
	return t, nil
}

func ScanSql(rows *sql.Rows, dest interface{}, logger zerolog.Logger) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return ErrScanDestinationPtr
	}
	if destVal.IsNil() {
		return ErrScanDestinationNil
	}
	_, err := validateType(destVal.Type(), reflect.Slice)
	if err != nil {
		logger.Error().Err(err).Msgf("error validating type, epected %v but got - %v", reflect.Slice, destVal.Type().Kind())
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	unique := map[string]any{}
	for rows.Next() {
		// TODO: Implement MapScan for sqlx
		// "github.com/jmoiron/sqlx"
		// err := rows.MapScan(result)
		// if err != nil {
		//     log.Fatal(err)
		// }
		if rows.Err() != nil {
			return rows.Err()
		}
		// Create a slice of reflect.Value to hold the values of the struct
		values := make([]interface{}, len(columns))
		row := make(map[string]any)
		for i := range values {
			values[i] = &values[i]
		}
		rows.Scan(values...)
		for i, v := range values {
			if v != nil {
				vType := reflect.TypeOf(v)
				unique[columns[i]] = vType
				switch v.(type) {
				case int64:
					v = v.(int64)
				case float64:
					v = v.(float64)
				case float32:
					v = v.(float32)
				case int:
					v = v.(int)
				case int16:
					v = v.(int16)
				case int8:
					v = v.(int8)
				case int32:
					v = v.(int32)
				case []byte:
					v = string(v.([]byte))
				case bool:
					v = v.(bool)
				case uint:
					v = v.(uint)
				case uint16:
					v = v.(uint16)
				case uint32:
					v = v.(uint32)
				case uint8:
					v = v.(uint8)
				case uint64:
					v = v.(uint64)
				case []uint64:
					v = v.([]uint64)
				case []uint16:
					v = v.([]uint16)
				case []uint32:
					v = v.([]uint32)
				case []uint:
					v = v.([]uint)
				case time.Time:
					v = v.(time.Time)
				case string:
					v = v.(string)
				case []string:
					v = v.([]string)
				default:
					if ok, s := dmutils.IsInt(v.(string)); ok {
						v = s
					} else if ok, s := dmutils.IsFloat(v.(string)); ok {
						v = s
					} else {
						bytes, err := json.Marshal(v)
						if err != nil {
							v = fmt.Sprintf("%v", v)
						}
						v = string(bytes)
					}
				}
			}

			row[columns[i]] = v
		}
		destVal.Elem().Set(reflect.Append(destVal.Elem(), reflect.ValueOf(row)))
	}
	return nil
}
