package mw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ginBodyLogger struct {
	// get all the methods implementation from the original one
	// override only the Write() method
	gin.ResponseWriter
	body bytes.Buffer
}

func (g *ginBodyLogger) Write(b []byte) (int, error) {
	g.body.Write(b)
	return g.ResponseWriter.Write(b)
}

func RequestLoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ginBodyLogger := &ginBodyLogger{
			body:           bytes.Buffer{},
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = ginBodyLogger
		var req interface{}
		if ctx.Request.Method != "GET" {
			if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
		}
		data, err := json.Marshal(req)
		if err != nil {
			panic(fmt.Errorf("err while marshaling req msg: %v", err))
		}
		ctx.Next()
		logger.WithFields(logrus.Fields{
			"status":       ctx.Writer.Status(),
			"method":       ctx.Request.Method,
			"path":         ctx.Request.URL.Path,
			"query_params": ctx.Request.URL.Query(),
			"req_body":     string(data),
			"res_body":     ginBodyLogger.body.String(),
		}).Info()
	}
}
