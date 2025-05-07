package datasvc

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const (
	VAPUS_ORGANIZATION_CATALOG_MAP = "vapus_ORGANIZATION_catalog_mapping"
)

// List of data collections in for of elastic search index names
const (
	VAPUS_ACCOUNT_INDEX                   = "vapus_accounts"
	VAPUS_DATASOURCE_INDEX                = "vapus_datasources"
	VAPUS_ORGANIZATION_INDEX              = "vapus_ORGANIZATIONs"
	VAPUS_DATASOURCE_METADATA_INDEX       = "vapus_datasource_metadata"
	VAPUS_ARTIFACTSOURCE_METADATA_INDEX   = "vapus_artifactsource_metadata"
	VAPUS_DATA_WORKER_INDEX               = "vapus_dataworkers"
	VAPUS_DATA_PRODUCTS_INDEX             = "vapus_dataproducts"
	VAPUS_USERS_INDEX                     = "vapus_users"
	VAPUS_DATA_CATALOGS_INDEX             = "vapus_datacatalogs"
	VAPUS_USER_GROUPS_INDEX               = "vapus_user_groups"
	VAPUS_DATAMARKETPLACE                 = "vapus_datamarketplace"
	VAPUS_USER_INVITES                    = "vapus_user_invites"
	VAPUS_JWT_INDEX                       = "vapus_jwt_log"
	VAPUS_DATA_WORKER_WORKFLOW_INDEX      = "vapus_dataworker_workflows"
	VAPUS_DATA_PRODUCT_DEPLOYMENT_INDEX   = "vapus_dataproduct_deployments"
	VAPUS_FEDERATED_CATALOG_INDEX         = "vapus_federated_catalog"
	DATA_PRODUCT_CONSUMPTION_METRICES     = "vapus_dataproduct_consumption_metrics"
	DATA_PRODUCT_CONSUMPTION_METRICES_RAW = "vapus_dataproduct_consumption_metrics_raw"
	VAPUS_DATA_WORKER_DEPLOYMENT_INDEX    = "vapus_dataworkers_deployments"
	VAPUS_DATA_CONTAINER_WORKER_INDEX     = "vapus_datacontainer_workers"
	VAPUS_AIMODEL_NODE_INDEX              = "vapus_aimodel_nodes"
)

// List of vaults for different resources
const (
	VAPUS_PLATFORM = "vapusdata-platform"
)

var INDEX_LIST = []string{
	VAPUS_ACCOUNT_INDEX,
	VAPUS_DATASOURCE_INDEX,
	VAPUS_ORGANIZATION_INDEX, VAPUS_DATASOURCE_METADATA_INDEX,
	VAPUS_DATA_WORKER_INDEX, VAPUS_DATA_PRODUCTS_INDEX,
	VAPUS_USERS_INDEX, VAPUS_DATA_CATALOGS_INDEX, VAPUS_USER_GROUPS_INDEX,
	VAPUS_DATAMARKETPLACE, VAPUS_USER_INVITES, VAPUS_JWT_INDEX, VAPUS_ARTIFACTSOURCE_METADATA_INDEX,
	VAPUS_DATA_WORKER_WORKFLOW_INDEX, VAPUS_DATA_PRODUCT_DEPLOYMENT_INDEX, VAPUS_FEDERATED_CATALOG_INDEX,
	DATA_PRODUCT_CONSUMPTION_METRICES, DATA_PRODUCT_CONSUMPTION_METRICES_RAW, VAPUS_AIMODEL_NODE_INDEX,
}
var denseDims = 512
var VAPUS_DATASOURCE_METADATA_INDEX_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataSourceSchemeVector": &types.DenseVectorProperty{
			Type: "dense_vector",
			Dims: &denseDims,
		},
		"dataSourceId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}
var VAPUS_ORGANIZATION_INDEX_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"OrganizationId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}

var VAPUS_AIMODEL_NODE_INDEX_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"modelNodeId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}

var DATA_PRODUCT_CONSUMPTION_METRICES_RAW_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataProductId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}

var VAPUS_DATASOURCE_INDEX_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataSourceId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}

var VAPUS_DATA_WORKER_INDEX_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataWorkerId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}

var VAPUS_DATA_PRODUCT_DEPLOYMENT_INDEX_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataProductId": &types.KeywordProperty{
			Type: "keyword",
		},
		"deploymentId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}

var VAPUS_FEDERATED_CATALOG_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataProductVector": &types.DenseVectorProperty{
			Type: "dense_vector",
			Dims: &denseDims,
		},
		"dataProductDescriptionVector": &types.DenseVectorProperty{
			Type: "dense_vector",
			Dims: &denseDims,
		},
	},
}
var VAPUS_ORGANIZATION_CATALOG_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataProductDescription": &types.DenseVectorProperty{
			Type: "dense_vector",
			Dims: &denseDims,
		},
		"dataProductId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}
var VAPUS_DATA_CATALOG_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"dataCatalogDescription": &types.DenseVectorProperty{
			Type: "dense_vector",
			Dims: &denseDims,
		},
		"dataCatalogId": &types.KeywordProperty{
			Type: "keyword",
		},
	},
}
var VAPUS_ARTIFACT_METADATA_INDEX_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"artifactsVector": &types.DenseVectorProperty{
			Type: "dense_vector",
			Dims: &denseDims,
		},
	},
}
var DATA_PRODUCT_CONSUMPTION_METRICES_MAPPING = &types.TypeMapping{
	Properties: map[string]types.Property{
		"consumptionMetrics": &types.DenseVectorProperty{
			Type: "dense_vector",
			Dims: &denseDims,
		},
	},
}
var INDEX_MAPPING = map[string]*types.TypeMapping{
	VAPUS_DATASOURCE_METADATA_INDEX:       VAPUS_DATASOURCE_METADATA_INDEX_MAPPING,
	VAPUS_ARTIFACTSOURCE_METADATA_INDEX:   VAPUS_ARTIFACT_METADATA_INDEX_MAPPING,
	VAPUS_FEDERATED_CATALOG_INDEX:         VAPUS_FEDERATED_CATALOG_MAPPING,
	VAPUS_ORGANIZATION_CATALOG_MAP:        VAPUS_ORGANIZATION_CATALOG_MAPPING,
	VAPUS_DATA_CATALOGS_INDEX:             VAPUS_DATA_CATALOG_MAPPING,
	DATA_PRODUCT_CONSUMPTION_METRICES:     DATA_PRODUCT_CONSUMPTION_METRICES_MAPPING,
	VAPUS_ORGANIZATION_INDEX:              VAPUS_ORGANIZATION_INDEX_MAPPING,
	VAPUS_DATA_WORKER_INDEX:               VAPUS_DATA_WORKER_INDEX_MAPPING,
	VAPUS_DATASOURCE_INDEX:                VAPUS_DATASOURCE_INDEX_MAPPING,
	DATA_PRODUCT_CONSUMPTION_METRICES_RAW: DATA_PRODUCT_CONSUMPTION_METRICES_RAW_MAPPING,
	VAPUS_AIMODEL_NODE_INDEX:              VAPUS_AIMODEL_NODE_INDEX_MAPPING,
}

const (
	SECRETENGINE = "secreteengine"
)
