package apppkgs

import (
	"fmt"
	"reflect"

	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

func GetAccountFilter(ctxClaim map[string]string, condition string) string {
	if ctxClaim != nil && ctxClaim[encryption.ClaimAccountKey] != "" {
		if condition == "" {
			return " owner_account = '" + ctxClaim[encryption.ClaimAccountKey] + "' "
		} else {
			return " owner_account = '" + ctxClaim[encryption.ClaimAccountKey] + "'" + " AND " + condition
		}
	}
	return condition
}

func GetByIdFilter(fieldId, val string, ctxClaim map[string]string) string {
	if fieldId == "" && val == "" {
		return GetAccountFilter(ctxClaim, "")
	}
	if fieldId == "" {
		fieldId = "vapus_id"
	}
	if ctxClaim != nil && ctxClaim[encryption.ClaimAccountKey] != "" {
		return fieldId + " = '" + val + "'" + " AND " + GetAccountFilter(ctxClaim, "")
	}
	return fieldId + " = '" + val + "'"
}

func VapusIdFilter() string {
	return "vapus_id = ?"
}

func GetOrderByQuery(field, order string) string {
	return fmt.Sprintf(" ORDER BY %s %s", field, order)
}

func BasePostFilterForamtting(condition string) string {
	if condition == "" {
		return GetOrderByQuery("created_at", "DESC")
	}
	return condition
}

func BuildSelectQuery(table, fields, condition, orderBy string) string {
	query := fmt.Sprintf("SELECT %s FROM %s", fields, table)
	if condition != "" {
		query += " WHERE " + condition
	}
	return query
}

func GetSQLConditons[T any](fieldId string, val T) string {
	if fieldId == "" && dmutils.IsNil(val) {
		return ""
	}

	result := ""
	switch reflect.TypeOf(val).Kind() {
	case reflect.String:
		result = fieldId + " = '" + fmt.Sprintf("%v", val) + "'"
	case reflect.Int:
		result = fieldId + " = '" + fmt.Sprintf("%d", val) + "'"
	case reflect.Int64:
		result = fieldId + " = '" + fmt.Sprintf("%d", val) + "'"
	case reflect.Float64:
		result = fieldId + " = '" + fmt.Sprintf("%f", val) + "'"
	case reflect.Bool:
		result = fieldId + " is " + fmt.Sprintf("%v", val)
	// case reflect.Slice:
	// 	v := reflect.ValueOf(val)
	// 	result = fieldId + " in " + fmt.Sprintf("[%v]", strings.Join(v, ","))
	default:
		result = fieldId + " = '" + fmt.Sprintf("%v", val) + "'"
	}
	return result
}
