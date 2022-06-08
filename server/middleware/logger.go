package middleware

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// realClock save request times
type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

func (rc *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

type LogOptions struct {
	Formatter      logrus.Formatter
	EnableStarting bool
}

// LoggingMiddleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type LoggingMiddleware struct {
	logger         *logrus.Logger
	clock          timer
	enableStarting bool
}

// NewLogger returns a new *LoggingMiddleware, yay!
func NewLogger(opts ...LogOptions) *LoggingMiddleware {
	var opt LogOptions
	if len(opts) == 0 {
		opt = LogOptions{}
	} else {
		opt = opts[0]
	}

	if opt.Formatter == nil {
		opt.Formatter = &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: time.RFC3339,
		}
	}

	log := logrus.New()
	//log.Formatter = new(logrus.JSONFormatter)
	//log.Formatter = new(logrus.TextFormatter) //default
	log.Formatter = opt.Formatter
	log.Level = logrus.TraceLevel
	//log.Out = os.Stdout

	return &LoggingMiddleware{
		logger:         log,
		clock:          &realClock{},
		enableStarting: opt.EnableStarting,
	}
}

// LoggingResponseWriter will encapsulate a standard ResponseWritter with a copy of its statusCode
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// ResponseWriterWrapper is supposed to capture statusCode from ResponseWriter
func ResponseWriterWrapper(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader is a surcharge of the ResponseWriter method
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	return lrw.ResponseWriter.Write(b)
}

// Logger Middleware implement mux middleware interface
func (m *LoggingMiddleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := logrus.NewEntry(m.logger)

		start := m.clock.Now()
		//token := r.Context().Value("token").(*models.Token)

		if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
			entry = entry.WithField("requestId", reqID)
		}

		requestLogger := entry.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.RequestURI,
		})

		fmt.Println()
		if m.enableStarting {
			entry.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.RequestURI,
				//"user":   token.UserId,
			}).Debugln("--------------------")
		}

		ctx := context.WithValue(r.Context(), "logger", requestLogger)
		r = r.WithContext(ctx)

		lw := ResponseWriterWrapper(w)

		next.ServeHTTP(lw, r)

		latency := time.Duration.Round(m.clock.Since(start), time.Millisecond)

		entry.WithFields(logrus.Fields{
			"status": lw.statusCode,
			"took":   latency,
		}).Debugln("--------------------")
	})
}

//// Logger is a gorilla/mux middleware to add log to the API
//func Logger(inner http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		//start := time.Now()
//		token := r.Context().Value("token").(*models.Token)
//
//		requestLogger := logrus.WithFields(logrus.Fields{"method": r.Method, "path": r.RequestURI, "status": wrapper.statusCode, "user_id": token.UserId})
//
//		ctx := context.WithValue(r.Context(), "logger", requestLogger)
//
//		r = r.WithContext(ctx)
//		inner.ServeHTTP(wrapper, r)
//		// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 286.219µs
//		//log.Printf("\"%s %s %s\" %d %d \"%s\" %s",
//		//	r.Method,
//		//	r.RequestURI,
//		//	r.Proto, // string "HTTP/1.1"
//		//	wrapper.statusCode,
//		//	r.ContentLength,
//		//	r.Header["User-Agent"],
//		//	time.Since(start),
//		//)
//
//	})
//}
