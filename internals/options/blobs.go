package options

type BlobOpsParams struct {
	BucketName       string
	Region           string
	ObjectName       string
	Data             []byte
	ObjectVersionId  string
	ByPassGovernance bool
}
