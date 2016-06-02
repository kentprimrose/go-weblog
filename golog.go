package weblog

import (
	"log"
	"net/http"
	"time"
)

type pipeWriter struct {
	writer http.ResponseWriter
	status int
	bytes  int
}

func (pw *pipeWriter) Status() int {
	return pw.status
}

func (pw *pipeWriter) Bytes() int {
	return pw.bytes
}

func (pw *pipeWriter) Header() http.Header {
	return pw.writer.Header()
}

func (pw *pipeWriter) Write(data []byte) (int, error) {
	rwBytes, rwErr := pw.writer.Write(data)
	pw.bytes += rwBytes
	return rwBytes, rwErr
}

func (pw *pipeWriter) WriteHeader(statusCode int) {
	pw.status = statusCode
	pw.writer.WriteHeader(statusCode)
}

func makePipeWriter(w http.ResponseWriter) pipeWriter {
	pw := pipeWriter{
		writer: w,
		status: 200,
	}
	return pw
}

// Handler for logging to specified logger.
func Handler(next http.Handler, l ...*log.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pw := makePipeWriter(w)
		t1 := time.Now()
		next.ServeHTTP(&pw, r)
		t2 := time.Now()
		if len(l) > 0 {
			l[0].Printf("%s %q %d %d %v\n",
				r.Method, r.URL.String(), pw.Status(), pw.Bytes(), t2.Sub(t1))
		} else {
			log.Printf("%s %q %d %d %v\n",
				r.Method, r.URL.String(), pw.Status(), pw.Bytes(), t2.Sub(t1))
		}
	}
	return http.HandlerFunc(fn)
}
