package create_db

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/eks"
// 	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
// 	"github.com/aws/aws-sdk-go-v2/service/rds"
// )

// // createEKSCluster creates a new EKS cluster. The user must supply:
// // - clusterName: name of the EKS cluster
// // - roleARN: the ARN of the IAM role that EKS will assume
// // - subnets: slice of subnet IDs for the cluster
// // - securityGroups: slice of security group IDs for the cluster
// func createEKSCluster(ctx context.Context, eksClient *eks.Client,
// 	clusterName, roleARN string,
// 	subnets, securityGroups []string,
// 	clusterVersion string,
// ) error {

// 	input := &eks.CreateClusterInput{
// 		Name:    aws.String(clusterName),
// 		RoleArn: aws.String(roleARN),
// 		ResourcesVpcConfig: &eksTypes.VpcConfigRequest{
// 			SubnetIds:        subnets,
// 			SecurityGroupIds: securityGroups,
// 		},
// 		// EKS version, e.g., "1.27"
// 		Version: aws.String(clusterVersion),
// 		// Optionally, you can add tags: map[string]string
// 	}

// 	out, err := eksClient.CreateCluster(ctx, input)
// 	if err != nil {
// 		return fmt.Errorf("failed to create EKS cluster: %w", err)
// 	}

// 	log.Printf("EKS cluster creation initiated: %s (status: %s)\n",
// 		aws.ToString(out.Cluster.Name),
// 		out.Cluster.Status,
// 	)
// 	return nil
// }

// // createRDSPostgresInstance creates a new RDS PostgreSQL instance.
// // - instanceID: a unique identifier for the DB instance
// // - dbName: the name of your database
// // - username/password: admin credentials
// // - dbSubnetGroupName: name of an existing DB subnet group
// func createRDSPostgresInstance(ctx context.Context, rdsClient *rds.Client,
// 	instanceID, dbName, username, password, dbSubnetGroupName string,
// ) error {

// 	input := &rds.CreateDBInstanceInput{
// 		DBInstanceIdentifier: aws.String(instanceID),
// 		DBName:               aws.String(dbName),
// 		Engine:               aws.String("postgres"),
// 		DBInstanceClass:      aws.String("db.t3.micro"), // example instance class
// 		AllocatedStorage:     aws.Int32(20),             // 20 GB
// 		MasterUsername:       aws.String(username),
// 		MasterUserPassword:   aws.String(password),

// 		// The subnet group you created or want to use
// 		DBSubnetGroupName: aws.String(dbSubnetGroupName),

// 		// Whether the instance has a public IP
// 		PubliclyAccessible: aws.Bool(false),

// 		// Additional configuration as needed...
// 		// BackupRetentionPeriod, StorageEncrypted, MultiAZ, etc.
// 	}

// 	out, err := rdsClient.CreateDBInstance(ctx, input)
// 	if err != nil {
// 		return fmt.Errorf("failed to create RDS instance: %w", err)
// 	}

// 	log.Printf("RDS Postgres creation initiated: %v (status: %v)\n",
// 		aws.ToString(out.DBInstance.DBInstanceIdentifier),
// 		out.DBInstance.DBInstanceStatus)
// 	return nil
// }

// func main() {
// 	ctx := context.Background()

// 	// 1. Load AWS configuration (region, credentials, etc.)
// 	cfg, err := config.LoadDefaultConfig(ctx,
// 		config.WithRegion("us-east-1"), // set your desired region
// 	)
// 	if err != nil {
// 		log.Fatalf("unable to load SDK config, %v", err)
// 	}

// 	// 2. Create EKS and RDS clients
// 	eksClient := eks.NewFromConfig(cfg)
// 	rdsClient := rds.NewFromConfig(cfg)

// 	// 3. Define parameters for EKS cluster creation
// 	eksClusterName := "my-eks-cluster"
// 	eksRoleARN := "arn:aws:iam::123456789012:role/EKSClusterRole"
// 	subnets := []string{"subnet-abc123", "subnet-def456"} // your VPC subnets
// 	secGroups := []string{"sg-abcdef1234567890"}          // security group(s)
// 	eksVersion := "1.27"

// 	if err := createEKSCluster(ctx, eksClient, eksClusterName, eksRoleARN, subnets, secGroups, eksVersion); err != nil {
// 		log.Fatalf("Error creating EKS cluster: %v", err)
// 	}
// 	log.Println("EKS cluster creation request submitted...")

// 	// 4. Define parameters for RDS Postgres creation
// 	rdsInstanceID := "my-postgres-db"
// 	dbName := "mydatabase"
// 	masterUsername := "admin"
// 	masterPassword := "SuperSecret123!"
// 	dbSubnetGroupName := "my-rds-subnet-group" // must exist or be created

// 	if err := createRDSPostgresInstance(ctx, rdsClient, rdsInstanceID, dbName,
// 		masterUsername, masterPassword, dbSubnetGroupName); err != nil {
// 		log.Fatalf("Error creating RDS instance: %v", err)
// 	}
// 	log.Println("RDS Postgres creation request submitted...")

// 	// 5. (Optional) You may want to poll for status or handle further steps,
// 	// such as creating EKS node groups, waiting for RDS to become available, etc.
// }

// type DBInstance struct {
// 	ActivityStreamEngineNativeAuditFieldsIncluded *bool
// 	ActivityStreamKinesisStreamName               *string
// 	ActivityStreamKmsKeyId                        *string
// 	ActivityStreamMode                            ActivityStreamMode
// 	ActivityStreamPolicyStatus                    ActivityStreamPolicyStatus
// 	ActivityStreamStatus                          ActivityStreamStatus
// 	AllocatedStorage                              *int32
// 	AssociatedRoles                               []DBInstanceRole
// 	AutoMinorVersionUpgrade                       *bool
// 	AutomaticRestartTime                          *time.Time
// 	AutomationMode                                AutomationMode
// 	AvailabilityZone                              *string
// 	AwsBackupRecoveryPointArn                     *string
// 	BackupRetentionPeriod                         *int32
// 	BackupTarget                                  *string
// 	CACertificateIdentifier                       *string
// 	CertificateDetails                            *CertificateDetails
// 	CharacterSetName                              *string
// 	CopyTagsToSnapshot                            *bool
// 	CustomIamInstanceProfile                      *string
// 	CustomerOwnedIpEnabled                        *bool
// 	DBClusterIdentifier                           *string
// 	DBInstanceArn                                 *string
// 	DBInstanceAutomatedBackupsReplications        []DBInstanceAutomatedBackupsReplication
// 	DBInstanceClass                               *string
// 	DBInstanceIdentifier                          *string
// 	DBInstanceStatus                              *string
// 	DBName                                        *string
// 	DBParameterGroups                             []DBParameterGroupStatus
// 	DBSecurityGroups                              []DBSecurityGroupMembership
// 	DBSubnetGroup                                 *DBSubnetGroup
// 	DBSystemId                                    *string
// 	DatabaseInsightsMode                          DatabaseInsightsMode
// 	DbInstancePort                                *int32
// 	DbiResourceId                                 *string
// 	DedicatedLogVolume                            *bool
// 	DeletionProtection                            *bool
// 	DomainMemberships                             []DomainMembership
// 	EnabledCloudwatchLogsExports                  []string
// 	Endpoint                                      *Endpoint
// 	Engine                                        *string
// 	EngineLifecycleSupport                        *string
// 	EngineVersion                                 *string
// 	EnhancedMonitoringResourceArn                 *string
// 	IAMDatabaseAuthenticationEnabled              *bool
// 	InstanceCreateTime                            *time.Time
// 	Iops                                          *int32
// 	IsStorageConfigUpgradeAvailable               *bool
// 	KmsKeyId                                      *string
// 	LatestRestorableTime                          *time.Time
// 	LicenseModel                                  *string
// 	ListenerEndpoint                              *Endpoint
// 	MasterUserSecret                              *MasterUserSecret
// 	MasterUsername                                *string
// 	MaxAllocatedStorage                           *int32
// 	MonitoringInterval                            *int32
// 	MonitoringRoleArn                             *string
// 	MultiAZ                                       *bool
// 	MultiTenant                                   *bool
// 	NcharCharacterSetName                         *string
// 	NetworkType                                   *string
// 	OptionGroupMemberships                        []OptionGroupMembership
// 	PendingModifiedValues                         *PendingModifiedValues
// 	PercentProgress                               *string
// 	PerformanceInsightsEnabled                    *bool
// 	PerformanceInsightsKMSKeyId                   *string
// 	PerformanceInsightsRetentionPeriod            *int32
// 	PreferredBackupWindow                         *string
// 	PreferredMaintenanceWindow                    *string
// 	ProcessorFeatures                             []ProcessorFeature
// 	PromotionTier                                 *int32
// 	PubliclyAccessible                            *bool
// 	ReadReplicaDBClusterIdentifiers               []string
// 	ReadReplicaDBInstanceIdentifiers              []string
// 	ReadReplicaSourceDBClusterIdentifier          *string
// 	ReadReplicaSourceDBInstanceIdentifier         *string
// 	ReplicaMode                                   ReplicaMode
// 	ResumeFullAutomationModeTime                  *time.Time
// 	SecondaryAvailabilityZone                     *string
// 	StatusInfos                                   []DBInstanceStatusInfo
// 	StorageEncrypted                              *bool
// 	StorageThroughput                             *int32
// 	StorageType                                   *string
// 	TagList                                       []Tag
// 	TdeCredentialArn                              *string
// 	Timezone                                      *string
// 	VpcSecurityGroups                             []VpcSecurityGroupMembership
// 	// noSmithyDocumentSerde
// }
