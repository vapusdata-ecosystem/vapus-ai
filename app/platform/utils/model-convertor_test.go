package utils

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	datamarketplace "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/datamarketplace"
// 	pb "github.com/vapusdata-ecosystem/vapusdata/platform/api/v1"
// )

// func TestpbtoOb(t *testing.T) {
// 	request := &pb.DataMarketplace{
// 		Name:        "test",
// 		DisplayName: "Test DataMarketplace",
// 	}

// 	expected := &datamarketplace.DataMarketplace{
// 		Name:        "test",
// 		DisplayName: "Test DataMarketplace",
// 	}

// 	result := pbtoOb(request)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodesPbtoOb(t *testing.T) {
// 	request := &pb.DataMarketplaceNode{
// 		Name:     "test",
// 		:      "http://example.com",
// 		Protocol: "http",
// 		Port:     8080,
// 		Attributes: &pb.DataStorageAttributes{
// 			StorageGoal:        pb.DataStorageGoal_DATA_STORAGE_GOAL_UNSPECIFIED,
// 			StorageServiceType: pb.DataStorageServicesTypes_DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED,
// 			ServiceName:        pb.DataStorageServices_DATA_STORAGE_SERVICES_UNSPECIFIED,
// 			ServiceProvider:    pb.DataStorageServiceProvider_DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED,
// 		},
// 		NodeCreds: []*pb.DataMarketplaceNodeCreds{},
// 	}

// 	expected := &datamarketplace.DMNode{
// 		Name:     "test",
// 		:      "http://example.com",
// 		Protocol: "http",
// 		Port:     8080,
// 		Attributes: &datamarketplace.DMNodeAttrtibutes{
// 			DMNodeType: "DATA_STORAGE_NODE_TYPE_UNSPECIFIED",
// 			DMNodeGoal: "DATA_STORAGE_GOAL_UNSPECIFIED",
// 			DMNodeSP: &datamarketplace.DMNodeSP{
// 				SvcName:     "DATA_STORAGE_SERVICES_UNSPECIFIED",
// 				SvcType:     "DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED",
// 				SvcProvider: "DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED",
// 			},
// 		},
// 		Credentials: []*datamarketplace.VDMNodeCreds{},
// 	}

// 	result := DMNodesPbtoOb(request)
// 	assert.Equal(t, expected, result)
// }

// func TestDMObToPb(t *testing.T) {
// 	dmObj := &datamarketplace.DataMarketplace{
// 		Name:        "test",
// 		DisplayName: "Test DataMarketplace",
// 	}

// 	dmNodes := []*datamarketplace.DMNode{
// 		{
// 			Name:   "test",
// 			:    "http://example.com",
// 			Port:   8080,
// 			Status: "STATUS_UNSPECIFIED",
// 			NodeId: "",
// 			Attributes: &datamarketplace.DMNodeAttrtibutes{
// 				DMNodeType: "DATA_STORAGE_NODE_TYPE_UNSPECIFIED",
// 				DMNodeGoal: "DATA_STORAGE_GOAL_UNSPECIFIED",
// 				DMNodeSP: &datamarketplace.DMNodeSP{
// 					SvcName:     "DATA_STORAGE_SERVICES_UNSPECIFIED",
// 					SvcType:     "DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED",
// 					SvcProvider: "DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED",
// 				},
// 			},
// 			Credentials: []*datamarketplace.VDMNodeCreds{},
// 		},
// 	}

// 	expected := &pb.DataMarketplace{
// 		Name:        "test",
// 		DisplayName: "Test DataMarketplace",
// 		MarketplaceId:      "",
// 		Nodes: []*pb.DataMarketplaceNode{
// 			{
// 				Name:     "test",
// 				:      "http://example.com",
// 				Port:     8080,
// 				Status:   "STATUS_UNSPECIFIED",
// 				NodeId:   "",
// 				NodeType: pb.DataMarketplaceNodeType_DATA_STORAGE_NODE_TYPE_UNSPECIFIED,
// 				Attributes: &pb.DataStorageAttributes{
// 					StorageGoal:        pb.DataStorageGoal_DATA_STORAGE_GOAL_UNSPECIFIED,
// 					StorageServiceType: pb.DataStorageServicesTypes_DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED,
// 					ServiceName:        pb.DataStorageServices_DATA_STORAGE_SERVICES_UNSPECIFIED,
// 					ServiceProvider:    pb.DataStorageServiceProvider_DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED,
// 				},
// 				NodeCreds: []*pb.DataMarketplaceNodeCreds{},
// 			},
// 		},
// 		Status: "STATUS_UNSPECIFIED",
// 	}

// 	result := DMObToPb(dmObj, dmNodes)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodesObToPb(t *testing.T) {
// 	dmnObj := &datamarketplace.DMNode{
// 		Name:   "test",
// 		:    "http://example.com",
// 		Port:   8080,
// 		Status: "STATUS_UNSPECIFIED",
// 		NodeId: "",
// 		Attributes: &datamarketplace.DMNodeAttrtibutes{
// 			DMNodeType: "DATA_STORAGE_NODE_TYPE_UNSPECIFIED",
// 			DMNodeGoal: "DATA_STORAGE_GOAL_UNSPECIFIED",
// 			DMNodeSP: &datamarketplace.DMNodeSP{
// 				SvcName:     "DATA_STORAGE_SERVICES_UNSPECIFIED",
// 				SvcType:     "DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED",
// 				SvcProvider: "DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED",
// 			},
// 		},
// 		Credentials: []*datamarketplace.VDMNodeCreds{},
// 	}

// 	expected := &pb.DataMarketplaceNode{
// 		Name:     "test",
// 		:      "http://example.com",
// 		Port:     8080,
// 		Status:   "STATUS_UNSPECIFIED",
// 		NodeId:   "",
// 		NodeType: pb.DataMarketplaceNodeType_JUST_STORAGE,
// 		Attributes: &pb.DataStorageAttributes{
// 			StorageGoal:        pb.DataStorageGoal_APPLICATION_DATA,
// 			StorageServiceType: pb.DataStorageServicesTypes_DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED,
// 			ServiceName:        pb.DataStorageServices_DATA_STORAGE_SERVICES_UNSPECIFIED,
// 			ServiceProvider:    pb.DataStorageServiceProvider_DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED,
// 		},
// 		NodeCreds: []*pb.DataMarketplaceNodeCreds{},
// 	}

// 	result := DMNodesObToPb(dmnObj)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodeCredPbToOb(t *testing.T) {
// 	request := &pb.DataMarketplaceNodeCreds{
// 		Name:              "test",
// 		CredentialVEngine: "engine",
// 		CredentialVPath:   "path",
// 		AccessScope:       pb.DMNodeCredAccessScope_DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED,
// 	}

// 	expected := &datamarketplace.VDMNodeCreds{
// 		Name:     "test",
// 		Priority: 0,
// 		UserType: "DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED",
// 		Secretpath: &datamarketplace.VDMNodeCredSecrets{
// 			SecretVPath:   "path",
// 			SecretVEngine: "engine",
// 		},
// 	}

// 	result := DMNodeCredPbToOb(request)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodeCredObToPb(t *testing.T) {
// 	obj := &datamarketplace.VDMNodeCreds{
// 		Name:     "test",
// 		Priority: 0,
// 		UserType: "DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED",
// 		Secretpath: &datamarketplace.VDMNodeCredSecrets{
// 			SecretVPath:   "path",
// 			SecretVEngine: "engine",
// 		},
// 	}

// 	expected := &pb.DataMarketplaceNodeCreds{
// 		Name:              "test",
// 		CredentialVEngine: "engine",
// 		CredentialVPath:   "path",
// 		AccessScope:       pb.DMNodeCredAccessScope_DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED,
// 	}

// 	result := DMNodeCredObToPb(obj)
// 	assert.Equal(t, expected, result)
// }
