// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/client/v1/client.proto

package clientv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v11 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/client/v1"
	v1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/common/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// GeneralsServiceName is the fully-qualified name of the GeneralsService service.
	GeneralsServiceName = "proto.common.v1.GeneralsService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GeneralsServiceResetProcedure is the fully-qualified name of the GeneralsService's Reset RPC.
	GeneralsServiceResetProcedure = "/proto.common.v1.GeneralsService/Reset"
	// GeneralsServiceReceiveCommandProcedure is the fully-qualified name of the GeneralsService's
	// ReceiveCommand RPC.
	GeneralsServiceReceiveCommandProcedure = "/proto.common.v1.GeneralsService/ReceiveCommand"
)

// GeneralsServiceClient is a client for the proto.common.v1.GeneralsService service.
type GeneralsServiceClient interface {
	// Reset resets the state
	Reset(context.Context, *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.EmptyResponse], error)
	// ReceiveCommand from the generals or commander
	ReceiveCommand(context.Context, *connect.Request[v11.ReceiveCommandRequest]) (*connect.Response[v1.EmptyResponse], error)
}

// NewGeneralsServiceClient constructs a client for the proto.common.v1.GeneralsService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGeneralsServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) GeneralsServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	generalsServiceMethods := v11.File_proto_client_v1_client_proto.Services().ByName("GeneralsService").Methods()
	return &generalsServiceClient{
		reset: connect.NewClient[v1.EmptyRequest, v1.EmptyResponse](
			httpClient,
			baseURL+GeneralsServiceResetProcedure,
			connect.WithSchema(generalsServiceMethods.ByName("Reset")),
			connect.WithClientOptions(opts...),
		),
		receiveCommand: connect.NewClient[v11.ReceiveCommandRequest, v1.EmptyResponse](
			httpClient,
			baseURL+GeneralsServiceReceiveCommandProcedure,
			connect.WithSchema(generalsServiceMethods.ByName("ReceiveCommand")),
			connect.WithClientOptions(opts...),
		),
	}
}

// generalsServiceClient implements GeneralsServiceClient.
type generalsServiceClient struct {
	reset          *connect.Client[v1.EmptyRequest, v1.EmptyResponse]
	receiveCommand *connect.Client[v11.ReceiveCommandRequest, v1.EmptyResponse]
}

// Reset calls proto.common.v1.GeneralsService.Reset.
func (c *generalsServiceClient) Reset(ctx context.Context, req *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.EmptyResponse], error) {
	return c.reset.CallUnary(ctx, req)
}

// ReceiveCommand calls proto.common.v1.GeneralsService.ReceiveCommand.
func (c *generalsServiceClient) ReceiveCommand(ctx context.Context, req *connect.Request[v11.ReceiveCommandRequest]) (*connect.Response[v1.EmptyResponse], error) {
	return c.receiveCommand.CallUnary(ctx, req)
}

// GeneralsServiceHandler is an implementation of the proto.common.v1.GeneralsService service.
type GeneralsServiceHandler interface {
	// Reset resets the state
	Reset(context.Context, *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.EmptyResponse], error)
	// ReceiveCommand from the generals or commander
	ReceiveCommand(context.Context, *connect.Request[v11.ReceiveCommandRequest]) (*connect.Response[v1.EmptyResponse], error)
}

// NewGeneralsServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGeneralsServiceHandler(svc GeneralsServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	generalsServiceMethods := v11.File_proto_client_v1_client_proto.Services().ByName("GeneralsService").Methods()
	generalsServiceResetHandler := connect.NewUnaryHandler(
		GeneralsServiceResetProcedure,
		svc.Reset,
		connect.WithSchema(generalsServiceMethods.ByName("Reset")),
		connect.WithHandlerOptions(opts...),
	)
	generalsServiceReceiveCommandHandler := connect.NewUnaryHandler(
		GeneralsServiceReceiveCommandProcedure,
		svc.ReceiveCommand,
		connect.WithSchema(generalsServiceMethods.ByName("ReceiveCommand")),
		connect.WithHandlerOptions(opts...),
	)
	return "/proto.common.v1.GeneralsService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GeneralsServiceResetProcedure:
			generalsServiceResetHandler.ServeHTTP(w, r)
		case GeneralsServiceReceiveCommandProcedure:
			generalsServiceReceiveCommandHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGeneralsServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGeneralsServiceHandler struct{}

func (UnimplementedGeneralsServiceHandler) Reset(context.Context, *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.EmptyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.common.v1.GeneralsService.Reset is not implemented"))
}

func (UnimplementedGeneralsServiceHandler) ReceiveCommand(context.Context, *connect.Request[v11.ReceiveCommandRequest]) (*connect.Response[v1.EmptyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.common.v1.GeneralsService.ReceiveCommand is not implemented"))
}
