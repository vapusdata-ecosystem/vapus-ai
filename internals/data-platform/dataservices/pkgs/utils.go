package pkgs

import "fmt"

func GetVapusQueryServerUri(dsId, schema, table string) string {
	return fmt.Sprintf("%s.%s.%s", dsId, schema, table)
}

func GetGeneralUri(schema, table string) string {
	return fmt.Sprintf("%s.%s", schema, table)
}
