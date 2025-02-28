// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: demo/v1/api.proto

package apiconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	// ApiServiceName is the fully-qualified name of the ApiService service.
	ApiServiceName = "demo.v1.ApiService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ApiServiceListClassesProcedure is the fully-qualified name of the ApiService's ListClasses RPC.
	ApiServiceListClassesProcedure = "/demo.v1.ApiService/ListClasses"
	// ApiServiceListRacesProcedure is the fully-qualified name of the ApiService's ListRaces RPC.
	ApiServiceListRacesProcedure = "/demo.v1.ApiService/ListRaces"
	// ApiServiceCreatePlayerProcedure is the fully-qualified name of the ApiService's CreatePlayer RPC.
	ApiServiceCreatePlayerProcedure = "/demo.v1.ApiService/CreatePlayer"
	// ApiServiceUpdatePlayerProcedure is the fully-qualified name of the ApiService's UpdatePlayer RPC.
	ApiServiceUpdatePlayerProcedure = "/demo.v1.ApiService/UpdatePlayer"
	// ApiServiceDeletePlayerProcedure is the fully-qualified name of the ApiService's DeletePlayer RPC.
	ApiServiceDeletePlayerProcedure = "/demo.v1.ApiService/DeletePlayer"
	// ApiServiceGetPlayerProcedure is the fully-qualified name of the ApiService's GetPlayer RPC.
	ApiServiceGetPlayerProcedure = "/demo.v1.ApiService/GetPlayer"
	// ApiServiceListPlayersAndCharactersProcedure is the fully-qualified name of the ApiService's
	// ListPlayersAndCharacters RPC.
	ApiServiceListPlayersAndCharactersProcedure = "/demo.v1.ApiService/ListPlayersAndCharacters"
	// ApiServiceCreateCharacterProcedure is the fully-qualified name of the ApiService's
	// CreateCharacter RPC.
	ApiServiceCreateCharacterProcedure = "/demo.v1.ApiService/CreateCharacter"
	// ApiServiceUpdateCharacterProcedure is the fully-qualified name of the ApiService's
	// UpdateCharacter RPC.
	ApiServiceUpdateCharacterProcedure = "/demo.v1.ApiService/UpdateCharacter"
	// ApiServiceDeleteCharacterProcedure is the fully-qualified name of the ApiService's
	// DeleteCharacter RPC.
	ApiServiceDeleteCharacterProcedure = "/demo.v1.ApiService/DeleteCharacter"
	// ApiServiceGetCharacterProcedure is the fully-qualified name of the ApiService's GetCharacter RPC.
	ApiServiceGetCharacterProcedure = "/demo.v1.ApiService/GetCharacter"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	apiServiceServiceDescriptor                        = v1.File_demo_v1_api_proto.Services().ByName("ApiService")
	apiServiceListClassesMethodDescriptor              = apiServiceServiceDescriptor.Methods().ByName("ListClasses")
	apiServiceListRacesMethodDescriptor                = apiServiceServiceDescriptor.Methods().ByName("ListRaces")
	apiServiceCreatePlayerMethodDescriptor             = apiServiceServiceDescriptor.Methods().ByName("CreatePlayer")
	apiServiceUpdatePlayerMethodDescriptor             = apiServiceServiceDescriptor.Methods().ByName("UpdatePlayer")
	apiServiceDeletePlayerMethodDescriptor             = apiServiceServiceDescriptor.Methods().ByName("DeletePlayer")
	apiServiceGetPlayerMethodDescriptor                = apiServiceServiceDescriptor.Methods().ByName("GetPlayer")
	apiServiceListPlayersAndCharactersMethodDescriptor = apiServiceServiceDescriptor.Methods().ByName("ListPlayersAndCharacters")
	apiServiceCreateCharacterMethodDescriptor          = apiServiceServiceDescriptor.Methods().ByName("CreateCharacter")
	apiServiceUpdateCharacterMethodDescriptor          = apiServiceServiceDescriptor.Methods().ByName("UpdateCharacter")
	apiServiceDeleteCharacterMethodDescriptor          = apiServiceServiceDescriptor.Methods().ByName("DeleteCharacter")
	apiServiceGetCharacterMethodDescriptor             = apiServiceServiceDescriptor.Methods().ByName("GetCharacter")
)

// ApiServiceClient is a client for the demo.v1.ApiService service.
type ApiServiceClient interface {
	ListClasses(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListClassesResponse], error)
	ListRaces(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListRacesResponse], error)
	CreatePlayer(context.Context, *connect.Request[v1.CreatePlayerRequest]) (*connect.Response[v1.Player], error)
	UpdatePlayer(context.Context, *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.Player], error)
	DeletePlayer(context.Context, *connect.Request[v1.DeletePlayerRequest]) (*connect.Response[emptypb.Empty], error)
	GetPlayer(context.Context, *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.Player], error)
	ListPlayersAndCharacters(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListPlayersAndCharactersResponse], error)
	CreateCharacter(context.Context, *connect.Request[v1.CreateCharacterRequest]) (*connect.Response[v1.Character], error)
	UpdateCharacter(context.Context, *connect.Request[v1.UpdateCharacterRequest]) (*connect.Response[v1.Character], error)
	DeleteCharacter(context.Context, *connect.Request[v1.DeleteCharacterRequest]) (*connect.Response[emptypb.Empty], error)
	GetCharacter(context.Context, *connect.Request[v1.GetCharacterRequest]) (*connect.Response[v1.Character], error)
}

// NewApiServiceClient constructs a client for the demo.v1.ApiService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewApiServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ApiServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &apiServiceClient{
		listClasses: connect.NewClient[emptypb.Empty, v1.ListClassesResponse](
			httpClient,
			baseURL+ApiServiceListClassesProcedure,
			connect.WithSchema(apiServiceListClassesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listRaces: connect.NewClient[emptypb.Empty, v1.ListRacesResponse](
			httpClient,
			baseURL+ApiServiceListRacesProcedure,
			connect.WithSchema(apiServiceListRacesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createPlayer: connect.NewClient[v1.CreatePlayerRequest, v1.Player](
			httpClient,
			baseURL+ApiServiceCreatePlayerProcedure,
			connect.WithSchema(apiServiceCreatePlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updatePlayer: connect.NewClient[v1.UpdatePlayerRequest, v1.Player](
			httpClient,
			baseURL+ApiServiceUpdatePlayerProcedure,
			connect.WithSchema(apiServiceUpdatePlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deletePlayer: connect.NewClient[v1.DeletePlayerRequest, emptypb.Empty](
			httpClient,
			baseURL+ApiServiceDeletePlayerProcedure,
			connect.WithSchema(apiServiceDeletePlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getPlayer: connect.NewClient[v1.GetPlayerRequest, v1.Player](
			httpClient,
			baseURL+ApiServiceGetPlayerProcedure,
			connect.WithSchema(apiServiceGetPlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listPlayersAndCharacters: connect.NewClient[emptypb.Empty, v1.ListPlayersAndCharactersResponse](
			httpClient,
			baseURL+ApiServiceListPlayersAndCharactersProcedure,
			connect.WithSchema(apiServiceListPlayersAndCharactersMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createCharacter: connect.NewClient[v1.CreateCharacterRequest, v1.Character](
			httpClient,
			baseURL+ApiServiceCreateCharacterProcedure,
			connect.WithSchema(apiServiceCreateCharacterMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateCharacter: connect.NewClient[v1.UpdateCharacterRequest, v1.Character](
			httpClient,
			baseURL+ApiServiceUpdateCharacterProcedure,
			connect.WithSchema(apiServiceUpdateCharacterMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteCharacter: connect.NewClient[v1.DeleteCharacterRequest, emptypb.Empty](
			httpClient,
			baseURL+ApiServiceDeleteCharacterProcedure,
			connect.WithSchema(apiServiceDeleteCharacterMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getCharacter: connect.NewClient[v1.GetCharacterRequest, v1.Character](
			httpClient,
			baseURL+ApiServiceGetCharacterProcedure,
			connect.WithSchema(apiServiceGetCharacterMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// apiServiceClient implements ApiServiceClient.
type apiServiceClient struct {
	listClasses              *connect.Client[emptypb.Empty, v1.ListClassesResponse]
	listRaces                *connect.Client[emptypb.Empty, v1.ListRacesResponse]
	createPlayer             *connect.Client[v1.CreatePlayerRequest, v1.Player]
	updatePlayer             *connect.Client[v1.UpdatePlayerRequest, v1.Player]
	deletePlayer             *connect.Client[v1.DeletePlayerRequest, emptypb.Empty]
	getPlayer                *connect.Client[v1.GetPlayerRequest, v1.Player]
	listPlayersAndCharacters *connect.Client[emptypb.Empty, v1.ListPlayersAndCharactersResponse]
	createCharacter          *connect.Client[v1.CreateCharacterRequest, v1.Character]
	updateCharacter          *connect.Client[v1.UpdateCharacterRequest, v1.Character]
	deleteCharacter          *connect.Client[v1.DeleteCharacterRequest, emptypb.Empty]
	getCharacter             *connect.Client[v1.GetCharacterRequest, v1.Character]
}

// ListClasses calls demo.v1.ApiService.ListClasses.
func (c *apiServiceClient) ListClasses(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListClassesResponse], error) {
	return c.listClasses.CallUnary(ctx, req)
}

// ListRaces calls demo.v1.ApiService.ListRaces.
func (c *apiServiceClient) ListRaces(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListRacesResponse], error) {
	return c.listRaces.CallUnary(ctx, req)
}

// CreatePlayer calls demo.v1.ApiService.CreatePlayer.
func (c *apiServiceClient) CreatePlayer(ctx context.Context, req *connect.Request[v1.CreatePlayerRequest]) (*connect.Response[v1.Player], error) {
	return c.createPlayer.CallUnary(ctx, req)
}

// UpdatePlayer calls demo.v1.ApiService.UpdatePlayer.
func (c *apiServiceClient) UpdatePlayer(ctx context.Context, req *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.Player], error) {
	return c.updatePlayer.CallUnary(ctx, req)
}

// DeletePlayer calls demo.v1.ApiService.DeletePlayer.
func (c *apiServiceClient) DeletePlayer(ctx context.Context, req *connect.Request[v1.DeletePlayerRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.deletePlayer.CallUnary(ctx, req)
}

// GetPlayer calls demo.v1.ApiService.GetPlayer.
func (c *apiServiceClient) GetPlayer(ctx context.Context, req *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.Player], error) {
	return c.getPlayer.CallUnary(ctx, req)
}

// ListPlayersAndCharacters calls demo.v1.ApiService.ListPlayersAndCharacters.
func (c *apiServiceClient) ListPlayersAndCharacters(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListPlayersAndCharactersResponse], error) {
	return c.listPlayersAndCharacters.CallUnary(ctx, req)
}

// CreateCharacter calls demo.v1.ApiService.CreateCharacter.
func (c *apiServiceClient) CreateCharacter(ctx context.Context, req *connect.Request[v1.CreateCharacterRequest]) (*connect.Response[v1.Character], error) {
	return c.createCharacter.CallUnary(ctx, req)
}

// UpdateCharacter calls demo.v1.ApiService.UpdateCharacter.
func (c *apiServiceClient) UpdateCharacter(ctx context.Context, req *connect.Request[v1.UpdateCharacterRequest]) (*connect.Response[v1.Character], error) {
	return c.updateCharacter.CallUnary(ctx, req)
}

// DeleteCharacter calls demo.v1.ApiService.DeleteCharacter.
func (c *apiServiceClient) DeleteCharacter(ctx context.Context, req *connect.Request[v1.DeleteCharacterRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.deleteCharacter.CallUnary(ctx, req)
}

// GetCharacter calls demo.v1.ApiService.GetCharacter.
func (c *apiServiceClient) GetCharacter(ctx context.Context, req *connect.Request[v1.GetCharacterRequest]) (*connect.Response[v1.Character], error) {
	return c.getCharacter.CallUnary(ctx, req)
}

// ApiServiceHandler is an implementation of the demo.v1.ApiService service.
type ApiServiceHandler interface {
	ListClasses(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListClassesResponse], error)
	ListRaces(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListRacesResponse], error)
	CreatePlayer(context.Context, *connect.Request[v1.CreatePlayerRequest]) (*connect.Response[v1.Player], error)
	UpdatePlayer(context.Context, *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.Player], error)
	DeletePlayer(context.Context, *connect.Request[v1.DeletePlayerRequest]) (*connect.Response[emptypb.Empty], error)
	GetPlayer(context.Context, *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.Player], error)
	ListPlayersAndCharacters(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListPlayersAndCharactersResponse], error)
	CreateCharacter(context.Context, *connect.Request[v1.CreateCharacterRequest]) (*connect.Response[v1.Character], error)
	UpdateCharacter(context.Context, *connect.Request[v1.UpdateCharacterRequest]) (*connect.Response[v1.Character], error)
	DeleteCharacter(context.Context, *connect.Request[v1.DeleteCharacterRequest]) (*connect.Response[emptypb.Empty], error)
	GetCharacter(context.Context, *connect.Request[v1.GetCharacterRequest]) (*connect.Response[v1.Character], error)
}

// NewApiServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewApiServiceHandler(svc ApiServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	apiServiceListClassesHandler := connect.NewUnaryHandler(
		ApiServiceListClassesProcedure,
		svc.ListClasses,
		connect.WithSchema(apiServiceListClassesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceListRacesHandler := connect.NewUnaryHandler(
		ApiServiceListRacesProcedure,
		svc.ListRaces,
		connect.WithSchema(apiServiceListRacesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceCreatePlayerHandler := connect.NewUnaryHandler(
		ApiServiceCreatePlayerProcedure,
		svc.CreatePlayer,
		connect.WithSchema(apiServiceCreatePlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceUpdatePlayerHandler := connect.NewUnaryHandler(
		ApiServiceUpdatePlayerProcedure,
		svc.UpdatePlayer,
		connect.WithSchema(apiServiceUpdatePlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceDeletePlayerHandler := connect.NewUnaryHandler(
		ApiServiceDeletePlayerProcedure,
		svc.DeletePlayer,
		connect.WithSchema(apiServiceDeletePlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceGetPlayerHandler := connect.NewUnaryHandler(
		ApiServiceGetPlayerProcedure,
		svc.GetPlayer,
		connect.WithSchema(apiServiceGetPlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceListPlayersAndCharactersHandler := connect.NewUnaryHandler(
		ApiServiceListPlayersAndCharactersProcedure,
		svc.ListPlayersAndCharacters,
		connect.WithSchema(apiServiceListPlayersAndCharactersMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceCreateCharacterHandler := connect.NewUnaryHandler(
		ApiServiceCreateCharacterProcedure,
		svc.CreateCharacter,
		connect.WithSchema(apiServiceCreateCharacterMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceUpdateCharacterHandler := connect.NewUnaryHandler(
		ApiServiceUpdateCharacterProcedure,
		svc.UpdateCharacter,
		connect.WithSchema(apiServiceUpdateCharacterMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceDeleteCharacterHandler := connect.NewUnaryHandler(
		ApiServiceDeleteCharacterProcedure,
		svc.DeleteCharacter,
		connect.WithSchema(apiServiceDeleteCharacterMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	apiServiceGetCharacterHandler := connect.NewUnaryHandler(
		ApiServiceGetCharacterProcedure,
		svc.GetCharacter,
		connect.WithSchema(apiServiceGetCharacterMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/demo.v1.ApiService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ApiServiceListClassesProcedure:
			apiServiceListClassesHandler.ServeHTTP(w, r)
		case ApiServiceListRacesProcedure:
			apiServiceListRacesHandler.ServeHTTP(w, r)
		case ApiServiceCreatePlayerProcedure:
			apiServiceCreatePlayerHandler.ServeHTTP(w, r)
		case ApiServiceUpdatePlayerProcedure:
			apiServiceUpdatePlayerHandler.ServeHTTP(w, r)
		case ApiServiceDeletePlayerProcedure:
			apiServiceDeletePlayerHandler.ServeHTTP(w, r)
		case ApiServiceGetPlayerProcedure:
			apiServiceGetPlayerHandler.ServeHTTP(w, r)
		case ApiServiceListPlayersAndCharactersProcedure:
			apiServiceListPlayersAndCharactersHandler.ServeHTTP(w, r)
		case ApiServiceCreateCharacterProcedure:
			apiServiceCreateCharacterHandler.ServeHTTP(w, r)
		case ApiServiceUpdateCharacterProcedure:
			apiServiceUpdateCharacterHandler.ServeHTTP(w, r)
		case ApiServiceDeleteCharacterProcedure:
			apiServiceDeleteCharacterHandler.ServeHTTP(w, r)
		case ApiServiceGetCharacterProcedure:
			apiServiceGetCharacterHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedApiServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedApiServiceHandler struct{}

func (UnimplementedApiServiceHandler) ListClasses(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListClassesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.ListClasses is not implemented"))
}

func (UnimplementedApiServiceHandler) ListRaces(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListRacesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.ListRaces is not implemented"))
}

func (UnimplementedApiServiceHandler) CreatePlayer(context.Context, *connect.Request[v1.CreatePlayerRequest]) (*connect.Response[v1.Player], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.CreatePlayer is not implemented"))
}

func (UnimplementedApiServiceHandler) UpdatePlayer(context.Context, *connect.Request[v1.UpdatePlayerRequest]) (*connect.Response[v1.Player], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.UpdatePlayer is not implemented"))
}

func (UnimplementedApiServiceHandler) DeletePlayer(context.Context, *connect.Request[v1.DeletePlayerRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.DeletePlayer is not implemented"))
}

func (UnimplementedApiServiceHandler) GetPlayer(context.Context, *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.Player], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.GetPlayer is not implemented"))
}

func (UnimplementedApiServiceHandler) ListPlayersAndCharacters(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.ListPlayersAndCharactersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.ListPlayersAndCharacters is not implemented"))
}

func (UnimplementedApiServiceHandler) CreateCharacter(context.Context, *connect.Request[v1.CreateCharacterRequest]) (*connect.Response[v1.Character], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.CreateCharacter is not implemented"))
}

func (UnimplementedApiServiceHandler) UpdateCharacter(context.Context, *connect.Request[v1.UpdateCharacterRequest]) (*connect.Response[v1.Character], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.UpdateCharacter is not implemented"))
}

func (UnimplementedApiServiceHandler) DeleteCharacter(context.Context, *connect.Request[v1.DeleteCharacterRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.DeleteCharacter is not implemented"))
}

func (UnimplementedApiServiceHandler) GetCharacter(context.Context, *connect.Request[v1.GetCharacterRequest]) (*connect.Response[v1.Character], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("demo.v1.ApiService.GetCharacter is not implemented"))
}
