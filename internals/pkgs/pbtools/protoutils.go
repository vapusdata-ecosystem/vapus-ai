package pbtools

import (
	"context"
	"fmt"
	"log"

	"buf.build/go/protoyaml"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"
)

var ErrorFinal = "errorFinal"

var LocalAuthzTokenCtx = "authorization:Bearer"

// HandleGrpcError is a utility function to handle grpc errors
func HandleGrpcError(err error, code grpccodes.Code) error {
	e, ok := err.(dmerrors.Error)
	if !ok {
		return grpcstatus.Error(code, err.Error())
	}
	return grpcstatus.Error(code, e.Error())
}

func HandleCtxCustomMessage(ctx context.Context, msgType string, msg ...string) context.Context {
	if ctx.Value(types.CUSTOM_MESSAGE) == nil {
		cm := map[string]any{msgType: msg}
		return context.WithValue(ctx, types.CUSTOM_MESSAGE, cm)
	}
	cm := ctx.Value(types.CUSTOM_MESSAGE).(map[string]any)
	cm[msgType] = append(cm[msgType].([]string), msg...)
	return context.WithValue(ctx, types.CUSTOM_MESSAGE, cm)
}

func GetRpcAuthFromCtx(ctx context.Context) (string, context.Context, error) {
	token, err := rpcauth.AuthFromMD(ctx, "bearer")
	if err != nil || token == "" {
		return "", ctx, grpcstatus.Error(grpccodes.Unauthenticated, "Authentication bearer token not found in request metadata")
	}
	return token, context.WithValue(ctx, LocalAuthzTokenCtx, token), nil
}

// HandleResponse is a utility function to handle the base response
func HandleDMResponse(ctx context.Context, opts ...string) *mpb.DMResponse {
	switch len(opts) {
	case 0:
		return &mpb.DMResponse{}
	case 1:
		return &mpb.DMResponse{
			Message: "Internal Server Error",
		}
	case 2:
		if ctx.Value(types.CUSTOM_MESSAGE) == nil {
			return &mpb.DMResponse{
				Message:      opts[0],
				DmStatusCode: opts[1],
			}
		}
		cm := ctx.Value(types.CUSTOM_MESSAGE).(map[string]any)
		return &mpb.DMResponse{
			Message:      opts[0],
			DmStatusCode: opts[1],
			CustomMessage: func(cms map[string]any) []*mpb.MapList {
				var cm []*mpb.MapList
				for k, v := range cms {
					cm = append(cm, &mpb.MapList{
						Key:    k,
						Values: v.([]string),
					})
				}
				return cm
			}(cm),
		}
	default:
		return &mpb.DMResponse{}
	}

}

func GetSvcDns(svcName string, namespace string, port int64) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", svcName, namespace, port)
}

// var ProtoYamlMarshaller = protojson.MarshalOptions{}
// var ProtoYamlUnMarshaller = protojson.MarshalOptions{}

func SwapNewContextWithAuthToken(ctx context.Context) context.Context {
	token, err := rpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		token = ""
	}
	log.Println("Token: ", token)
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+token))
}

func SwapHttpContextWithAuthToken(ctx context.Context, token string) context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
}

func NewInBgCtxWithAuthToken(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "))
	}
	fCtx := metadata.NewIncomingContext(context.Background(), md)

	val := ctx.Value(encryption.JwtCtxClaimKey)
	claims, ok := val.(map[string]string)
	if ok {
		log.Println("Claims: ", claims)
		fCtx = context.WithValue(fCtx, encryption.JwtCtxClaimKey, claims)
	}
	return fCtx
}

func NewInToDoCtxWithAuthToken(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.NewIncomingContext(context.TODO(), metadata.Pairs("authorization", "Bearer "))
	}
	fCtx := metadata.NewIncomingContext(context.TODO(), md)

	val := ctx.Value(encryption.JwtCtxClaimKey)
	claims, ok := val.(map[string]string)
	if ok {
		log.Println("Claims: ", claims)
		fCtx = context.WithValue(fCtx, encryption.JwtCtxClaimKey, claims)
	}
	return fCtx
}

func NewInCancelCtxWithAuthToken(ctx context.Context) (context.Context, context.CancelFunc) {
	nCtx, cancel := context.WithCancel(context.Background())
	log.Println("NewInCancelCtxWithAuthToken------------------------------------")
	log.Println("Ctx: ++++++++++++++++++++++++++++++++++++", ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Println("MD: ", md)
	log.Println("OK: ", ok)
	if !ok {
		return metadata.NewIncomingContext(nCtx, metadata.Pairs("authorization", "Bearer ")), cancel
	}
	fCtx := metadata.NewIncomingContext(nCtx, md)

	val := ctx.Value(encryption.JwtCtxClaimKey)
	claims, ok := val.(map[string]string)
	if ok {
		log.Println("Claims: ", claims)
		fCtx = context.WithValue(fCtx, encryption.JwtCtxClaimKey, claims)
	}
	return fCtx, cancel
}

func GetPbAnyToGoAny(pbAny *anypb.Any) any {
	var result any
	if pbAny == nil {
		return ""
	}
	listValue := &structpb.ListValue{}

	if err := pbAny.UnmarshalTo(listValue); err == nil {
		sliceResult := []any{}
		for _, value := range listValue.Values {
			switch val := value.Kind.(type) {
			case *structpb.Value_StringValue:
				sliceResult = append(sliceResult, val.StringValue)
			case *structpb.Value_NumberValue:
				sliceResult = append(sliceResult, val.NumberValue)
			default:
				log.Printf("Unknown type in list value: %v", value)
			}
		}
		return sliceResult
	}

	// Handle scalar types
	stringValue := &wrapperspb.StringValue{}
	intValue := &wrapperspb.Int64Value{}

	if err := pbAny.UnmarshalTo(stringValue); err == nil {
		result = stringValue.Value
		return result
	}

	if err := pbAny.UnmarshalTo(intValue); err == nil {
		result = intValue.Value
		return result
	}
	return result
}

var ProtoJsonMarshaller = protojson.MarshalOptions{
	Indent:          "",
	UseEnumNumbers:  false,
	EmitUnpopulated: true,
	Multiline:       false,
}

var ProtoJsonUnMarshaller = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

var ProtoYamlMarshaller = protoyaml.MarshalOptions{
	Indent:          2,
	UseEnumNumbers:  false,
	EmitUnpopulated: true,
}

var ProtoYamlUnMarshaller = protoyaml.UnmarshalOptions{
	DiscardUnknown: true,
}

func ProtoJsonMarshal(message protoreflect.ProtoMessage) ([]byte, error) {
	return ProtoJsonMarshaller.Marshal(message)
}

func ProtoJsonUnMarshal(bytes []byte, message protoreflect.ProtoMessage) error {
	return ProtoJsonUnMarshaller.Unmarshal(bytes, message)
}

func ProtoYamlMarshal(message protoreflect.ProtoMessage) ([]byte, error) {
	return ProtoYamlMarshaller.Marshal(message)
}

func ProtoYamlUnMarshal(bytes []byte, message protoreflect.ProtoMessage) error {
	return ProtoYamlUnMarshaller.Unmarshal(bytes, message)
}
