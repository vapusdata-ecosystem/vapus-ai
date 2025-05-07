package mistral

func getToolType(t string) string {
	val, ok := ToolTypeMap[t]
	if !ok {
		return ToolTypeFunction.String()
	}
	return val.String()
}
