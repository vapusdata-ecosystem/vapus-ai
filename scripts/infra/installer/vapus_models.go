package vapusmodel

type SetupPostgressInstance struct {
	InstanceID          string   `json:"instanceID" yaml:"instanceID" toml:"instanceID"`
	DBName              string   `json:"dbName" yaml:"dbName" toml:"dbName"`
	Username            string   `json:"username" yaml:"username" toml:"username"`
	Password            string   `json:"password" yaml:"password" toml:"password"`
	DBSubnetGroupName   string   `json:"dbSubnetGroupName" yaml:"dbSubnetGroupName" toml:"dbSubnetGroupName"`
	DBInstanceClass     string   `json:"dbInstanceClass" yaml:"dbInstanceClass" toml:"dbInstanceClass"`
	AllocatedStorage    int32    `json:"allocatedStorage" yaml:"allocatedStorage" toml:"allocatedStorage"`
	PubliclyAccessible  bool     `json:"publiclyAccessible" yaml:"publiclyAccessible" toml:"publiclyAccessible"`
	DBSecurityGroups    []string `json:"dBSecurityGroups,omitempty" yaml:"dBSecurityGroups,omitempty" toml:"dBSecurityGroups,omitempty"`
	DeletionProtection  bool     `json:"deletionProtection,omitempty" yaml:"deletionProtection,omitempty" toml:"deletionProtection,omitempty"`
	Port                int32    `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
	VpcSecurityGroupIds []string `json:"vpcSecurityGroupIds,omitempty" yaml:"vpcSecurityGroupIds,omitempty" toml:"vpcSecurityGroupIds,omitempty"`
	StorageEncrypted    bool     `json:"storageEncrypted,omitempty" yaml:"storageEncrypted,omitempty" toml:"storageEncrypted,omitempty"`
	// DBParameterGroupName string   `json:"dBParameterGroupName" yaml:"dBParameterGroupName" toml:"dBParameterGroupName"`
}

type AWSConfig struct {
	Region    string `json:"region" yaml:"region" toml:"region"`
	AccessKey string `json:"accessKey" yaml:"accessKey" toml:"accessKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey" toml:"secretKey"`
}

type SetupBucketInput struct {
	BucketName                 string
	Region                     string
	GrantFullControl           string
	GrantRead                  string
	GrantReadACP               string
	GrantWrite                 string
	GrantWriteACP              string
	ObjectLockEnabledForBucket bool
}

// type CreateDBInstanceInput struct {
// 	DBInstanceClass                    *string
// 	DBInstanceIdentifier               *string
// 	Engine                             *string
// 	AllocatedStorage                   *int32
// 	AutoMinorVersionUpgrade            *bool
// 	AvailabilityZone                   *string
// 	BackupRetentionPeriod              *int32
// 	BackupTarget                       *string
// 	CACertificateIdentifier            *string
// 	CharacterSetName                   *string
// 	CopyTagsToSnapshot                 *bool
// 	CustomIamInstanceProfile           *string
// 	DBClusterIdentifier                *string
// 	DBName                             *string
// 	DBParameterGroupName               *string
// 	DBSecurityGroups                   []string
// 	DBSubnetGroupName                  *string
// 	DBSystemId                         *string
// 	DatabaseInsightsMode               types.DatabaseInsightsMode
// 	DedicatedLogVolume                 *bool
// 	DeletionProtection                 *bool
// 	Domain                             *string
// 	DomainAuthSecretArn                *string
// 	DomainDnsIps                       []string
// 	DomainFqdn                         *string
// 	DomainIAMRoleName                  *string
// 	DomainOu                           *string
// 	EnableCloudwatchLogsExports        []string
// 	EnableCustomerOwnedIp              *bool
// 	EnableIAMDatabaseAuthentication    *bool
// 	EnablePerformanceInsights          *bool
// 	EngineLifecycleSupport             *string
// 	EngineVersion                      *string
// 	Iops                               *int32
// 	KmsKeyId                           *string
// 	LicenseModel                       *string
// 	ManageMasterUserPassword           *bool
// 	MasterUserPassword                 *string
// 	MasterUserSecretKmsKeyId           *string
// 	MasterUsername                     *string
// 	MaxAllocatedStorage                *int32
// 	MonitoringInterval                 *int32
// 	MonitoringRoleArn                  *string
// 	MultiAZ                            *bool
// 	MultiTenant                        *bool
// 	NcharCharacterSetName              *string
// 	NetworkType                        *string
// 	OptionGroupName                    *string
// 	PerformanceInsightsKMSKeyId        *string
// 	PerformanceInsightsRetentionPeriod *int32
// 	Port                               *int32
// 	PreferredBackupWindow              *string
// 	PreferredMaintenanceWindow         *string
// 	ProcessorFeatures                  []types.ProcessorFeature
// 	PromotionTier                      *int32
// 	PubliclyAccessible                 *bool
// 	StorageEncrypted                   *bool
// 	StorageThroughput                  *int32
// 	StorageType                        *string
// 	Tags                               []types.Tag
// 	TdeCredentialArn                   *string
// 	TdeCredentialPassword              *string
// 	Timezone                           *string
// 	VpcSecurityGroupIds                []string
// 	// noSmithyDocumentSerde
// }

// type CreateBucketInput struct {
// 	Bucket                     *string
// 	ACL                        types.BucketCannedACL
// 	CreateBucketConfiguration  *types.CreateBucketConfiguration
// 	GrantFullControl           *string
// 	GrantRead                  *string
// 	GrantReadACP               *string
// 	GrantWrite                 *string
// 	GrantWriteACP              *string
// 	ObjectLockEnabledForBucket *bool
// 	ObjectOwnership            types.ObjectOwnership
// 	// noSmithyDocumentSerde
// }
