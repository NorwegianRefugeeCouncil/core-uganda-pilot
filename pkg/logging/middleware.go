package logging

import (
	"bufio"
	"github.com/emicklei/go-restful/v3"
	"github.com/felixge/httpsnoop"
	"github.com/lithammer/shortuuid/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"time"
)

func UseOperationLogging() func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
		route := request.SelectedRoute()
		var operation string
		if route != nil {
			operation = route.Operation()
		} else {
			operation = ""
		}
		ctx := WithOperation(request.Request.Context(), operation)
		request.Request = request.Request.WithContext(ctx)
		chain.ProcessFilter(request, response)
	}
}

func UseRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		const keyRequestId = "X-Request-Id"
		reqId := req.Header.Get(keyRequestId)
		if len(reqId) == 0 {
			newId := shortuuid.New()
			req.Header.Set(keyRequestId, newId)
		}
		ctx := WithRequestID(req.Context(), req.Header.Get(keyRequestId))
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func UseRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		l := NewLogger(ctx)

		start := time.Now()
		logger, w := makeLogger(w)

		next.ServeHTTP(w, req)

		statusCode := logger.Status()

		fields := []zapcore.Field{
			zap.Int("status", statusCode),
			zap.String("latency", time.Since(start).String()),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
			zap.String("host", req.Host),
			zap.String("remote_address", req.RemoteAddr),
		}

		for key, strings := range req.Header {
			logKey := "header." + key
			if len(strings) == 1 {
				fields = append(fields, zap.String(logKey, strings[0]))
			} else {
				fields = append(fields, zap.Strings(logKey, strings))
			}
		}

		n := statusCode
		switch {
		case n >= 500:
			l.Error("server error", fields...)
		case n >= 400:
			l.Error("client error", fields...)
		case n >= 300:
			l.Info("redirection", fields...)
		default:
			l.Info("success", fields...)
		}

	})
}

func makeLogger(w http.ResponseWriter) (*responseLogger, http.ResponseWriter) {
	logger := &responseLogger{w: w, status: http.StatusOK}
	return logger, httpsnoop.Wrap(w, httpsnoop.Hooks{
		Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return logger.Write
		},
		WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return logger.WriteHeader
		},
	})
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, rw, err := l.w.(http.Hijacker).Hijack()
	if err == nil && l.status == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and
		// WriteHeader has not been called yet
		l.status = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}
