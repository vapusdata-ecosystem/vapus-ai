package salesforce

type ListObjects struct {
	Sobjects []struct {
		Name string `json:"name"`
	} `json:"sobjects"`
}

type DescribeField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DescribeResponse holds the object schema
type DescribeResponse struct {
	Fields []DescribeField `json:"fields"`
}
