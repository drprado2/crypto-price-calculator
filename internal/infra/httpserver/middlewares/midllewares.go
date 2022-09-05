package middlewares

import (
	"context"
	"crypto-price-calculator/internal/ctxutils"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.Status = code
	rec.ResponseWriter.WriteHeader(code)
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				applog.Logger(r.Context()).Errorf("Panic occurs in path %v, error: %v, stacktrace: %s", r.RequestURI, err, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := context.WithValue(r.Context(), "httpMethod", r.Method)
		reqCtx = context.WithValue(reqCtx, "httpPath", r.RequestURI)

		logger := applog.Logger(reqCtx)

		logger.
			WithFields(logrus.Fields{"method": r.Method, "path": r.RequestURI}).
			Infof("Handling request, method: %v, path: %v", r.Method, r.RequestURI)

		writter := &StatusRecorder{
			w, http.StatusOK,
		}
		next.ServeHTTP(writter, r.WithContext(reqCtx))

		if writter.Status >= 500 {
			logger.WithFields(logrus.Fields{"httpStatusCode": writter.Status, "requestSuccess": false}).Errorf("request fineshed with errors, status code: %v", writter.Status)
		} else if writter.Status >= 400 {
			logger.WithFields(logrus.Fields{"httpStatusCode": writter.Status, "requestSuccess": false}).Warnf("request fineshed with warnings, status code: %v", writter.Status)
		} else {
			logger.WithFields(logrus.Fields{"httpStatusCode": writter.Status, "requestSuccess": true}).Infof("request fineshed with success, status code: %v", writter.Status)
		}
	})
}

func CidMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get("x-cid")
		if cid == "" {
			cid = uuid.New().String()
		}
		reqCtx := ctxutils.WithCid(r.Context(), cid)
		next.ServeHTTP(w, r.WithContext(reqCtx))
	})
}

func SpanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := apptracer.StartOperation(r.Context(), fmt.Sprintf("%s::%s", r.Method, r.RequestURI), apptracer.SpanKindInternal)
		defer span.Finish()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
