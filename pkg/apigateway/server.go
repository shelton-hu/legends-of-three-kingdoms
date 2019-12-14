package apigateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/shelton-hu/util/hashutil"
	"github.com/shelton-hu/util/iputil"
	"github.com/shelton-hu/util/pointerutil"

	staticSpec "github.com/shelton-hu/legends-of-three-kingdoms/pkg/apigateway/spec"
	staticSwaggerUI "github.com/shelton-hu/legends-of-three-kingdoms/pkg/apigateway/swagger-ui"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/config"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/constants"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/gerr"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/manager"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pb"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pi"
	prisma "github.com/shelton-hu/legends-of-three-kingdoms/pkg/prisma/mysql-prisma-client"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/util/senderutil"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/version"
)

type Server struct {
	config.IAMConfig
}

type register struct {
	f        func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	endpoint string
}

func Serve() {
	ctx := context.Background()
	version.PrintVersionInfo(func(s string, i ...interface{}) {
		logger.Info(ctx, s, i...)
	})

	logger.Info(ctx, "Api service start http://%s:%d", constants.ApiGatewayHost, constants.ApiGatewayPort)
	logger.Info(ctx, "IAM manager start http://%s:%d", constants.IAMManagerHost, constants.IAMManagerPort)
	logger.Info(ctx, "Room manager start http://%s:%d", constants.RoomManagerHost, constants.RoomManagerPort)
	logger.Info(ctx, "Game manager start http://%s:%d", constants.GameManagerHost, constants.GameManagerPort)

	s := Server{}

	if err := s.run(); err != nil {
		logger.Critical(ctx, "Api gateway run failed: %+v", err)
		panic(err)
	}
}

const (
	Authorization = "Authorization"
	RequestIdKey  = "X-Request-Id"
	xForwardedFor = "X-Forwarded-For"
)

func log() gin.HandlerFunc {
	ctx := context.Background()
	l := logger.NewLogger()
	l.HideCallstack()
	return func(c *gin.Context) {
		requestID := uuid.New()
		c.Request.Header.Set(RequestIdKey, requestID)
		c.Writer.Header().Set(RequestIdKey, requestID)

		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		logStr := fmt.Sprintf("%s | %3d | %v | %s | %s %s %s",
			requestID,
			statusCode,
			latency,
			clientIP, method,
			path,
			c.Errors.String(),
		)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			l.Warn(ctx, logStr)
		case statusCode >= 500:
			l.Error(ctx, logStr)
		default:
			l.Info(ctx, logStr)
		}
	}
}

func serveMuxSetSender(mux *runtime.ServeMux, key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// 1.skip some special route
		switch req.URL.Path {
		case "/v1/oauth2/token":
			mux.ServeHTTP(w, req)
			return
		}

		// 2.get & validate sender
		ctx := req.Context()
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		auth := strings.SplitN(req.Header.Get(Authorization), " ", 2)
		if auth[0] != "Bearer" {
			err := gerr.New(ctx, gerr.Unauthenticated, gerr.ErrorAuthFailure)
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		sender, err := senderutil.Validate(key, auth[1])
		if err != nil {
			if err == senderutil.ErrExpired {
				err = gerr.New(ctx, gerr.Unauthenticated, gerr.ErrorAccessTokenExpired)
			} else {
				err = gerr.New(ctx, gerr.Unauthenticated, gerr.ErrorAuthFailure)
			}
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		user, err := pi.Global().MysqlPrisma(ctx).User(prisma.UserWhereUniqueInput{
			ID: &sender.UserId,
		}).Exec(ctx)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		if hashutil.GetStringMd5(auth[1]) != pointerutil.GetString(user.Token) {
			err = gerr.New(ctx, gerr.Unauthenticated, gerr.ErrorAuthFailure)
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		// 3.set & delete header's key
		req.Header.Set(senderutil.SenderKey, sender.ToJson())
		req.Header.Del(Authorization)

		// 4.route
		mux.ServeHTTP(w, req)
	})
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Critical(context.Background(), "Panic recovered: %+v\n%s", err, string(httprequest))
				c.JSON(500, gin.H{
					"title": "Error",
					"err":   err,
				})
			}
		}()
		c.Next() // execute all the handlers
	}
}

func handleSwagger() http.Handler {
	ns := vfs.NameSpace{}
	ns.Bind("/", mapfs.New(staticSwaggerUI.Files), "/", vfs.BindReplace)
	ns.Bind("/", mapfs.New(staticSpec.Files), "/", vfs.BindBefore)
	return http.StripPrefix("/swagger-ui", http.FileServer(httpfs.New(ns)))
}

func (s *Server) run() error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(log())
	r.Use(recovery())
	r.Any("/swagger-ui/*filepath", gin.WrapH(handleSwagger()))
	r.Any("/v1/*filepath", gin.WrapH(s.mainHandler()))

	return r.Run(fmt.Sprintf(":%d", constants.ApiGatewayPort))
}

func (s *Server) mainHandler() http.Handler {
	ctx := context.Background()
	var gwmux = runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			ip := iputil.GetRequestRemoteIP(req)
			return metadata.Pairs(
				senderutil.SenderKey, req.Header.Get(senderutil.SenderKey),
				RequestIdKey, req.Header.Get(RequestIdKey),
				xForwardedFor, ip,
			)
		}),
	)
	var opts = manager.ClientOptions
	var err error

	for _, r := range []register{
		{
			pb.RegisterTokenServiceHandlerFromEndpoint,
			fmt.Sprintf("%s:%d", constants.IAMManagerHost, constants.IAMManagerPort),
		},
	} {
		err = r.f(ctx, gwmux, r.endpoint, opts)
		if err != nil {
			err = errors.WithStack(err)
			logger.Error(ctx, "Dial [%s] failed: %+v", r.endpoint, err)
		}
	}

	mux := http.NewServeMux()

	mux.Handle("/", serveMuxSetSender(gwmux, s.IAMConfig.SecretKey))

	return formWrapper(mux)
}

// Ref: https://github.com/grpc-ecosystem/grpc-gateway/issues/7#issuecomment-358569373
func formWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			jsonMap := make(map[string]interface{}, len(r.Form))
			for k, v := range r.Form {
				if len(v) > 0 {
					jsonMap[k] = v[0]
				}
			}
			jsonBody, err := json.Marshal(jsonMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			r.Body = ioutil.NopCloser(bytes.NewReader(jsonBody))
			r.ContentLength = int64(len(jsonBody))
			r.Header.Set("Content-Type", "application/json")
		}

		h.ServeHTTP(w, r)
	})
}
