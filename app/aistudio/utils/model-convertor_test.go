package utils

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	datamesh "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/datamesh"
// 	pb "github.com/vapusdata-ecosystem/vapusdata/aistudio/api/v1"
// )

// func TestpbtoOb(t *testing.T) {
// 	request := &pb.DataMesh{
// 		Name:        "test",
// 		DisplayName: "Test DataMesh",
// 	}

// 	expected := &datamesh.DataMesh{
// 		Name:        "test",
// 		DisplayName: "Test DataMesh",
// 	}

// 	result := pbtoOb(request)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodesPbtoOb(t *testing.T) {
// 	request := &pb.DataMeshNode{
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
// 		NodeCreds: []*pb.DataMeshNodeCreds{},
// 	}

// 	expected := &datamesh.DMNode{
// 		Name:     "test",
// 		:      "http://example.com",
// 		Protocol: "http",
// 		Port:     8080,
// 		Attributes: &datamesh.DMNodeAttrtibutes{
// 			DMNodeType: "DATA_STORAGE_NODE_TYPE_UNSPECIFIED",
// 			DMNodeGoal: "DATA_STORAGE_GOAL_UNSPECIFIED",
// 			DMNodeSP: &datamesh.DMNodeSP{
// 				SvcName:     "DATA_STORAGE_SERVICES_UNSPECIFIED",
// 				SvcType:     "DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED",
// 				SvcProvider: "DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED",
// 			},
// 		},
// 		Credentials: []*datamesh.VDMNodeCreds{},
// 	}

// 	result := DMNodesPbtoOb(request)
// 	assert.Equal(t, expected, result)
// }

// func TestDMObToPb(t *testing.T) {
// 	dmObj := &datamesh.DataMesh{
// 		Name:        "test",
// 		DisplayName: "Test DataMesh",
// 	}

// 	dmNodes := []*datamesh.DMNode{
// 		{
// 			Name:   "test",
// 			:    "http://example.com",
// 			Port:   8080,
// 			Status: "STATUS_UNSPECIFIED",
// 			NodeId: "",
// 			Attributes: &datamesh.DMNodeAttrtibutes{
// 				DMNodeType: "DATA_STORAGE_NODE_TYPE_UNSPECIFIED",
// 				DMNodeGoal: "DATA_STORAGE_GOAL_UNSPECIFIED",
// 				DMNodeSP: &datamesh.DMNodeSP{
// 					SvcName:     "DATA_STORAGE_SERVICES_UNSPECIFIED",
// 					SvcType:     "DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED",
// 					SvcProvider: "DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED",
// 				},
// 			},
// 			Credentials: []*datamesh.VDMNodeCreds{},
// 		},
// 	}

// 	expected := &pb.DataMesh{
// 		Name:        "test",
// 		DisplayName: "Test DataMesh",
// 		MeshId:      "",
// 		Nodes: []*pb.DataMeshNode{
// 			{
// 				Name:     "test",
// 				:      "http://example.com",
// 				Port:     8080,
// 				Status:   "STATUS_UNSPECIFIED",
// 				NodeId:   "",
// 				NodeType: pb.DataMeshNodeType_DATA_STORAGE_NODE_TYPE_UNSPECIFIED,
// 				Attributes: &pb.DataStorageAttributes{
// 					StorageGoal:        pb.DataStorageGoal_DATA_STORAGE_GOAL_UNSPECIFIED,
// 					StorageServiceType: pb.DataStorageServicesTypes_DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED,
// 					ServiceName:        pb.DataStorageServices_DATA_STORAGE_SERVICES_UNSPECIFIED,
// 					ServiceProvider:    pb.DataStorageServiceProvider_DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED,
// 				},
// 				NodeCreds: []*pb.DataMeshNodeCreds{},
// 			},
// 		},
// 		Status: "STATUS_UNSPECIFIED",
// 	}

// 	result := DMObToPb(dmObj, dmNodes)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodesObToPb(t *testing.T) {
// 	dmnObj := &datamesh.DMNode{
// 		Name:   "test",
// 		:    "http://example.com",
// 		Port:   8080,
// 		Status: "STATUS_UNSPECIFIED",
// 		NodeId: "",
// 		Attributes: &datamesh.DMNodeAttrtibutes{
// 			DMNodeType: "DATA_STORAGE_NODE_TYPE_UNSPECIFIED",
// 			DMNodeGoal: "DATA_STORAGE_GOAL_UNSPECIFIED",
// 			DMNodeSP: &datamesh.DMNodeSP{
// 				SvcName:     "DATA_STORAGE_SERVICES_UNSPECIFIED",
// 				SvcType:     "DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED",
// 				SvcProvider: "DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED",
// 			},
// 		},
// 		Credentials: []*datamesh.VDMNodeCreds{},
// 	}

// 	expected := &pb.DataMeshNode{
// 		Name:     "test",
// 		:      "http://example.com",
// 		Port:     8080,
// 		Status:   "STATUS_UNSPECIFIED",
// 		NodeId:   "",
// 		NodeType: pb.DataMeshNodeType_JUST_STORAGE,
// 		Attributes: &pb.DataStorageAttributes{
// 			StorageGoal:        pb.DataStorageGoal_APPLICATION_DATA,
// 			StorageServiceType: pb.DataStorageServicesTypes_DATA_STORAGE_SERVICES_TYPES_UNSPECIFIED,
// 			ServiceName:        pb.DataStorageServices_DATA_STORAGE_SERVICES_UNSPECIFIED,
// 			ServiceProvider:    pb.DataStorageServiceProvider_DATA_STORAGE_SERVICE_PROVIDER_UNSPECIFIED,
// 		},
// 		NodeCreds: []*pb.DataMeshNodeCreds{},
// 	}

// 	result := DMNodesObToPb(dmnObj)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodeCredPbToOb(t *testing.T) {
// 	request := &pb.DataMeshNodeCreds{
// 		Name:              "test",
// 		CredentialVEngine: "engine",
// 		CredentialVPath:   "path",
// 		AccessScope:       pb.DMNodeCredAccessScope_DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED,
// 	}

// 	expected := &datamesh.VDMNodeCreds{
// 		Name:     "test",
// 		Priority: 0,
// 		UserType: "DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED",
// 		Secretpath: &datamesh.VDMNodeCredSecrets{
// 			SecretVPath:   "path",
// 			SecretVEngine: "engine",
// 		},
// 	}

// 	result := DMNodeCredPbToOb(request)
// 	assert.Equal(t, expected, result)
// }

// func TestDMNodeCredObToPb(t *testing.T) {
// 	obj := &datamesh.VDMNodeCreds{
// 		Name:     "test",
// 		Priority: 0,
// 		UserType: "DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED",
// 		Secretpath: &datamesh.VDMNodeCredSecrets{
// 			SecretVPath:   "path",
// 			SecretVEngine: "engine",
// 		},
// 	}

// 	expected := &pb.DataMeshNodeCreds{
// 		Name:              "test",
// 		CredentialVEngine: "engine",
// 		CredentialVPath:   "path",
// 		AccessScope:       pb.DMNodeCredAccessScope_DM_NODE_CRED_ACCESS_SCOPE_UNSPECIFIED,
// 	}

// 	result := DMNodeCredObToPb(obj)
// 	assert.Equal(t, expected, result)
// }
