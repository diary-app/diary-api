package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

func RequestLoggerMiddleware(log *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Debugf("Path: %s, request body: %s", c.Request.URL, string(body))
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseLoggerMiddleware(l *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		statusCode := c.Writer.Status()
		l.Debugf("Status: %v, response body: %v", statusCode, blw.body.String())
	}
}
