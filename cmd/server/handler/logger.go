package handler

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type loggerRW struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *loggerRW) Status() int {
	return rw.status
}

func (rw *loggerRW) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func loggerHandler(next http.HandlerFunc, logger *zerolog.Logger) http.HandlerFunc {
	lgr := logger.With().Str("section", "handler.logger").Logger()

	return func(rw http.ResponseWriter, r *http.Request) {

		start := time.Now()
		rrw := &loggerRW{ResponseWriter: rw}

		next.ServeHTTP(rrw, r)

		lgr.Info().
			Int("status", rrw.status).
			Str("method", r.Method).
			Str("path", r.URL.EscapedPath()).
			Dur("duration", time.Since(start)).
			Send()
	}
}
