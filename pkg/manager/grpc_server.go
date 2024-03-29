package manager

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/gerr"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/util/ctxutil"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/util/senderutil"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/version"
)

type checkerT func(ctx context.Context, req interface{}) error
type builderT func(ctx context.Context, req interface{}) interface{}

var (
	defaultChecker checkerT
	defaultBuilder builderT
)

type GrpcServer struct {
	ServiceName    string
	Port           int
	showErrorCause bool
	checker        checkerT
	builder        builderT
}

type RegisterCallBack func(*grpc.Server)

func NewGrpcServer(serviceName string, port int) *GrpcServer {
	return &GrpcServer{
		ServiceName:    serviceName,
		Port:           port,
		showErrorCause: false,
		checker:        defaultChecker,
		builder:        defaultBuilder,
	}
}

func (g *GrpcServer) ShowErrorCause(b bool) *GrpcServer {
	g.showErrorCause = b
	return g
}

func (g *GrpcServer) WithChecker(c checkerT) *GrpcServer {
	g.checker = c
	return g
}

func (g *GrpcServer) WithBuilder(b builderT) *GrpcServer {
	g.builder = b
	return g
}

func (g *GrpcServer) Serve(callback RegisterCallBack, opt ...grpc.ServerOption) {
	ctx := context.Background()
	version.PrintVersionInfo(func(s string, i ...interface{}) {
		logger.Info(ctx, s, i)
	})
	logger.Info(ctx, "Service [%s] start listen at port [%d]", g.ServiceName, g.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		err = errors.WithStack(err)
		logger.Critical(ctx, "failed to listen: %+v", err)
	}

	builtinOptions := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.MaxRecvMsgSize(100 * 1024 * 1024),
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			g.unaryServerLogInterceptor(),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if g.checker != nil {
					err = g.checker(ctx, req)
					if err != nil {
						return
					}
				}

				return handler(ctx, req)
			},
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if g.builder != nil {
					req = g.builder(ctx, req)
				}
				return handler(ctx, req)
			},
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Critical(ctx, "GRPC server recovery with error: %+v", p)
					logger.Critical(ctx, string(debug.Stack()))
					if e, ok := p.(error); ok {
						return gerr.NewWithDetail(ctx, gerr.Internal, e, gerr.ErrorInternalError)
					}
					return gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Critical(ctx, "GRPC server recovery with error: %+v", p)
					logger.Critical(ctx, string(debug.Stack()))
					if e, ok := p.(error); ok {
						return gerr.NewWithDetail(ctx, gerr.Internal, e, gerr.ErrorInternalError)
					}
					return gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
				}),
			),
		),
	}

	grpcServer := grpc.NewServer(append(opt, builtinOptions...)...)
	reflection.Register(grpcServer)
	callback(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Critical(ctx, "%+v", err)
	}
}

var (
	jsonPbMarshaller = &jsonpb.Marshaler{
		OrigName: true,
	}
)

func (g *GrpcServer) unaryServerLogInterceptor() grpc.UnaryServerInterceptor {
	showErrorCause := g.showErrorCause

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		s := senderutil.GetSenderFromContext(ctx)
		requestId := ctxutil.GetRequestId(ctx)
		forwardedFor := ctxutil.GetForwardedFor(ctx)
		ctx = ctxutil.SetRequestId(ctx, requestId)
		ctx = senderutil.ContextWithSender(ctx, s)
		ctx = ctxutil.SetForwardedFor(ctx, forwardedFor)

		method := strings.Split(info.FullMethod, "/")
		action := method[len(method)-1]
		if p, ok := req.(proto.Message); ok {
			if content, err := jsonPbMarshaller.MarshalToString(p); err != nil {
				logger.Error(ctx, "Failed to marshal proto message to string [%s] [%+v] [%+v]", action, s, err)
			} else {
				logger.Info(ctx, "Request received [%s] [%+v] [%s]", action, s, content)
			}
		}
		start := time.Now()

		resp, err := handler(ctx, req)

		elapsed := time.Since(start)
		logger.Info(ctx, "Handled request [%s] [%+v] exec_time is [%s]", action, s, elapsed)
		if e, ok := status.FromError(err); ok {
			if e.Code() != codes.OK {
				logger.Debug(ctx, "Response is error: %s, %s", e.Code().String(), e.Message())
				if !showErrorCause {
					err = gerr.ClearErrorCause(err)
				}
			}
		}
		return resp, err
	}
}
