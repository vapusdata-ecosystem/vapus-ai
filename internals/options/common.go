package options

type FieldMapping struct {
	FileId       string `json:"fileId"`
	CurrentField string `json:"currentField"`
	NewField     string `json:"newField"`
}
