package generic

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/boj/redistore"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"gopkg.in/matryer/try.v1"
)

type Server struct {
	options            options.ServerOptions
	name               string
	address            string
	listener           net.Listener
	GoRestfulContainer *restful.Container
	NonGoRestfulMux    *mux.Router
	handler            http.Handler
	sessionStore       sessions.Store
}

type Middleware func(next http.Handler) http.Handler

func NewGenericServer(options options.ServerOptions, name string) (*Server, error) {

	ctx := logging.WithServerName(context.Background(), name)
	l := logging.NewLogger(ctx)

	srv := &Server{
		name:    name,
		options: options,
	}

	if len(options.Secrets.Hash) != len(options.Secrets.Block) {
		return nil, fmt.Errorf("number of hash keys must be equal to number of block keys")
	}

	var keyPairs [][]byte
	for i := range options.Secrets.Hash {
		hashKey := options.Secrets.Hash[i]
		keyPairs = append(keyPairs, []byte(hashKey)[0:32])
		blockKey := options.Secrets.Block[i]
		keyPairs = append(keyPairs, []byte(blockKey)[0:32])
	}

	if options.Cache.Redis != nil {

		l.Info("configuring redis store")

		pool := &redis.Pool{
			MaxActive:   500,
			MaxIdle:     500,
			IdleTimeout: 5 * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				if err != nil {
					l.Error("failed to get redis connection", zap.Error(err))
				}
				return err
			},
			Dial: func() (redis.Conn, error) {
				var redisOptions []redis.DialOption
				if len(options.Cache.Redis.Password) > 0 {
					redisOptions = append(redisOptions, redis.DialPassword(options.Cache.Redis.Password))
				}
				return redis.Dial("tcp", options.Cache.Redis.Address, redisOptions...)
			},
		}

		conn := pool.Get()
		defer conn.Close()

		err := try.Do(func(attempt int) (retry bool, err error) {
			l.Debug("testing redis connection")
			_, err = conn.Do("PING")
			if err != nil {
				l.Error("redis connection test failed", zap.Error(err))
				time.Sleep(time.Duration(attempt) * time.Second)
				return attempt < 5, err
			}
			return false, nil
		})

		if err != nil {
			l.Error("failed to test redis connection. reached max attempts", zap.Error(err))
			panic(err)
		}

		redisStore, err := redistore.NewRediStoreWithPool(pool, keyPairs...)
		if err != nil {
			l.Error("failed to create redis store", zap.Error(err))
			panic(err)

		}
		if options.Cache.Redis.MaxLength != 0 {
			redisStore.SetMaxLength(options.Cache.Redis.MaxLength)
		}
		redisStore.Options.Secure = true
		redisStore.Options.HttpOnly = true
		redisStore.Options.SameSite = http.SameSiteNoneMode

		srv.sessionStore = redisStore

	} else {
		srv.sessionStore = sessions.NewCookieStore(keyPairs...)
	}

	var address string
	// If we have both host and port, use that as the address for the net listener
	if len(options.Host) != 0 && options.Port != 0 {
		address = fmt.Sprintf("%s:%d", options.Host, options.Port)
	} else if len(options.Host) != 0 {
		// if we have only an address, let the newListener give us a random port
		address = fmt.Sprintf("%s:0", options.Host)
	} else if options.Port != 0 {
		// if we have only a port, we listen to 0.0.0.0
		address = fmt.Sprintf(":%d", options.Port)
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		l.Error("failed to start listener", zap.Error(err), zap.String("address", address))
		return nil, err
	}

	// Retrieve the address from the listener. Perhaps the listener
	// listens on a random port (if options.Port = 0)
	srv.address = listener.Addr().String()
	srv.listener = listener

	container := restful.NewContainer()
	container.Filter(logging.UseOperationLogging())
	srv.GoRestfulContainer = container
	srv.NonGoRestfulMux = mux.NewRouter()
	srv.NonGoRestfulMux.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		utils.ErrorResponse(w, &meta.StatusError{
			ErrStatus: meta.Status{
				Status:  meta.StatusFailure,
				Message: "Not found",
				Reason:  meta.StatusReasonNotFound,
				Code:    http.StatusNotFound,
			},
		})
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:     options.Cors.AllowedOrigins,
		AllowCredentials:   options.Cors.AllowCredentials,
		AllowedMethods:     options.Cors.AllowedMethods,
		Debug:              options.Cors.Debug,
		OptionsPassthrough: options.Cors.OptionsPassthrough,
		AllowedHeaders:     options.Cors.AllowedHeaders,
		MaxAge:             options.Cors.MaxAge,
		ExposedHeaders:     options.Cors.ExposedHeaders,
	})
	corsHandler.Log = &CorsLogger{}

	middleware := chainMiddleware(
		func(next http.Handler) http.Handler {
			return logging.UseRequestLogging(next)
		},
		func(next http.Handler) http.Handler {
			return corsHandler.Handler(next)
		},
		func(next http.Handler) http.Handler {
			return handlers.CompressHandler(next)
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				req = req.WithContext(logging.WithServerName(req.Context(), name))
				next.ServeHTTP(w, req)
			})
		},
		func(next http.Handler) http.Handler {
			return logging.UseRequestID(next)
		},
		func(next http.Handler) http.Handler {
			return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(next)
		},
	)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		l := logging.NewLogger(ctx)
		path := req.URL.Path

		// check if go-restful can handle this request
		for _, ws := range container.RegisteredWebServices() {
			if strings.HasPrefix(path, ws.RootPath()) && (len(path) == len(ws.RootPath()) || path[len(ws.RootPath())] == '/') {
				l.Debug("satisfied by GoRestful with webservice",
					zap.String("path", path),
					zap.String("root_path", ws.RootPath()),
				)
				container.ServeHTTP(w, req)
				return
			}
		}

		// skip go-restful
		l.Debug("satisfied by NonGoRestful",
			zap.String("path", path),
		)
		srv.NonGoRestfulMux.ServeHTTP(w, req)

	}))
	srv.handler = handler

	return srv, nil

}

type CorsLogger struct {
}

func (c CorsLogger) Printf(s string, i ...interface{}) {
	logging.NewLogger(context.TODO()).Debug(fmt.Sprintf(s, i...), zap.String("middleware", "cors"))
}

var _ cors.Logger = CorsLogger{}

func (g Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.handler.ServeHTTP(w, req)
}

func (g Server) SessionStore() sessions.Store {
	return g.sessionStore
}

func (g Server) Address() string {
	return g.address
}

func (g Server) Port() int {
	parts := strings.Split(g.address, ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])
	return port
}

func (g Server) Start(ctx context.Context) {

	ctx = logging.WithServerName(ctx, g.name)
	l := logging.NewLogger(ctx).
		With(zap.String("address", g.listener.Addr().String())).
		With(zap.Bool("tls", g.options.TLS.Enabled))

	config := restfulspec.Config{
		WebServices: g.GoRestfulContainer.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/openapi.json",
		ModelTypeNameHandler: func(t reflect.Type) (string, bool) {
			return t.Name(), true
		},
	}
	g.GoRestfulContainer.Add(restfulspec.NewOpenAPIService(config))

	l.Info("starting server")

	go func() {
		var err error
		if g.options.TLS.Enabled {
			err = http.ServeTLS(g.listener, g.handler, g.options.TLS.Cert.Path, g.options.TLS.Key.Path)
		} else {
			err = http.Serve(g.listener, g.handler)
		}
		if err != nil {
			l.With(zap.Error(err)).Info("server stopped")
			if !errors.Is(err, net.ErrClosed) {
				panic(err)
			}
		}
	}()
	go func() {
		<-ctx.Done()
		if err := g.listener.Close(); err != nil {
			l.Error("failed to close listener", zap.Error(err))
			panic(err)
		}
	}()
}

// chainMiddleware provides syntactic sugar to create a new middleware
// which will be the result of chaining the ones received as parameters.
func chainMiddleware(mw ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last.ServeHTTP(w, r)
		})
	}
}
