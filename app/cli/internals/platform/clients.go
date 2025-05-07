package plclient

import (
	"fmt"
	"net/url"
	"time"

	"buf.build/go/protoyaml"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	appcl "github.com/vapusdata-ecosystem/vapusdata/core/app/grpcclients"
	gwcl "github.com/vapusdata-ecosystem/vapusdata/core/app/httpcls"
)

var AgentGoals = map[string][]interface{}{
	"account":        getAgentOps(pb.AccountAgentActions_name),
	"datacompliance": getAgentOps(nil),
	"authz":          getAgentOps(pb.AuthzAgentActions_name),
	"authorization":  getAgentOps(pb.AccessTokenAgentUtility_name),
	"utility":        getAgentOps(nil),
}

type ActionHandlerOpts struct {
	ParentCmd   string
	Args        []string
	Action      string
	File        string
	La          bool
	AccessToken string
	SearchQ     string
	Params      map[string]string
	Resource    string
}

type VapusCtlClient struct {
	Host               string
	ClientManager      *appcl.VapusSvcInternalClients
	CaCertFile         string
	ClientCertFile     string
	ClientKeyFile      string
	ValidTill          time.Time
	Error              error
	logger             zerolog.Logger
	ResourceActionMap  map[string][]interface{}
	inputFormat        string
	ActionHandler      ActionHandlerOpts
	protoyamlUnMarshal protoyaml.UnmarshalOptions
	protoyamlMarshal   protoyaml.MarshalOptions
	fileBytes          []byte
	GwClient           *gwcl.VapusHttpClient
	protojsonMarshal   protojson.MarshalOptions
}

func getAgentOps(enum_map map[int32]string) []interface{} {
	var ops []interface{}
	if enum_map == nil {
		return ops
	}
	for _, v := range enum_map {
		ops = append(ops, v)
	}
	return ops
}

func NewPlatFormClient(params map[string]string, logger zerolog.Logger) (*VapusCtlClient, error) {
	urlParsed, err := url.Parse(params["url"])
	if err != nil {
		return nil, err
	}
	// namespace, ok := params["namespace"]
	// if !ok {
	// 	return nil, errors.New("namespace is required, missing from the context")
	// }
	// port, ok := params["port"]
	// if !ok {
	// 	return nil, errors.New("port is required, missing from the context")
	// }
	// portI, err := strconv.Atoi(port)
	// if err != nil {
	// 	return nil, errors.Join(err, errors.New("port is not a valid integer"))
	// }
	dns := fmt.Sprintf("Connecting to vapusdata instance at - %s", urlParsed.String())
	// dns = "localhost:9013"

	// telnet, err := net.DialTimeout("tcp", dns, 1*time.Second)
	// if err != nil {
	// 	return nil, err
	// }
	// defer telnet.Close()

	// grpcClient := pbtools.NewGrpcClient(logger,
	// 	pbtools.ClientWithInsecure(true),
	// 	pbtools.ClientWithServiceAddress(dns))
	gwcls, err := gwcl.New(params["url"], logger)
	if err != nil {
		// return nil, err
		return nil, nil
	}

	cl := &VapusCtlClient{
		Host: dns,
		// PlConn:             pb.NewPlatformServiceClient(grpcClient.Connection),
		// platformGrpcClient: grpcClient,
		// UserConn:           pb.NewUserManagementServiceClient(grpcClient.Connection),
		// OrganizationConn:         dpb.NewOrganizationServiceClient(grpcClient.Connection),
		logger:             logger,
		protoyamlUnMarshal: protoyaml.UnmarshalOptions{},
		protoyamlMarshal: protoyaml.MarshalOptions{
			EmitUnpopulated: true,
		},
		protojsonMarshal: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
			UseEnumNumbers:  false,
		},
		ActionHandler: ActionHandlerOpts{},
		GwClient:      gwcls,
	}
	return cl, nil
}

// func (x *VapusCtlClient) setAIStudioClient(url string) error {
// 	svcinfo, err := x.PlConn.PlatformServicesInfo(context.Background(), &pb.PlatformServicesRequest{})
// 	if err != nil {
// 		return err
// 	}
// 	if svcinfo.GetNetworkParams() != nil || len(svcinfo.GetNetworkParams()) < 1 {
// 		for _, svc := range svcinfo.GetNetworkParams() {
// 			if svc.SvcTag == mpb.VapusSvcs_AISTUDIO {
// 				telnet, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", url, svc.Port), 1*time.Second)
// 				if err != nil {
// 					return err
// 				}
// 				defer telnet.Close()
// 				x.aiStudioGrpcClient = pbtools.NewGrpcClient(x.logger,
// 					pbtools.ClientWithInsecure(true),
// 					pbtools.ClientWithServiceAddress(fmt.Sprintf("%s:%d", url, svc.Port)))
// 				x.AIStudioConn = aipb.NewAIAgentStudioClient(x.aiStudioGrpcClient.Connection)
// 				return nil
// 			}
// 		}
// 	}
// 	return errors.New("AI Studio service not found")
// }

func (x *VapusCtlClient) Close() {
	return
}

func (x *VapusCtlClient) PrintDescribe(data protoreflect.ProtoMessage, resource string) {
	bytes, err := x.protoyamlMarshal.Marshal(data)
	if err != nil {
		x.logger.Error().Msgf("Error in marshaling %v details", resource)
	}
	x.logger.Info().Msgf("\n%s", string(bytes))
}
