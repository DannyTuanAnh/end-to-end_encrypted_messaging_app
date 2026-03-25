package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(data []byte) (n int, err error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func LoggerMiddleware() gin.HandlerFunc {

	logPath := "../../internal/logs/http.log"

	logger := zerolog.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5,    //days
		Compress:   true, // disable by default
		LocalTime:  true, // use local time for timestamps
	}).With().Timestamp().Logger()

	return func(ctx *gin.Context) {
		start := time.Now()
		contentType := ctx.GetHeader("Content-Type")
		requestBody := make(map[string]any)
		var formFiles []map[string]any

		// Content-Type: multipart/form-data
		if strings.HasPrefix(contentType, "multipart/form-data") {

			if err := ctx.Request.ParseMultipartForm(32 << 20); err == nil && ctx.Request.MultipartForm != nil {

				// Lấy giá trị text fields
				for key, vals := range ctx.Request.MultipartForm.Value {
					if len(vals) == 1 {
						requestBody[key] = vals[0]
					} else {
						requestBody[key] = vals
					}
				}

				// Lấy file from request
				for field, files := range ctx.Request.MultipartForm.File {
					for _, f := range files {
						fileInfo := map[string]any{
							"field":        field,
							"filename":     f.Filename,
							"size":         formatSize(f.Size),
							"content_type": f.Header.Get("Content-Type"),
						}

						formFiles = append(formFiles, fileInfo)
					}
				}

				if len(formFiles) > 0 {
					requestBody["form_files"] = formFiles
				}
			}
		} else {

			// Content-Type: application/x-www-form-urlencoded
			// Content-Type: application/json

			bodyBytes, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to read request body")
			}

			ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			if strings.HasPrefix(contentType, "application/json") {
				// Content-Type: application/json

				err = json.Unmarshal(bodyBytes, &requestBody)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to parse JSON body")
				}

			} else {
				// Content-Type: application/x-www-form-urlencoded

				values, err := url.ParseQuery(string(bodyBytes))
				if err != nil {
					logger.Error().Err(err).Msg("Failed to parse form data")
				}

				for key, vals := range values {
					if len(vals) == 1 {
						requestBody[key] = vals[0]
					} else {
						requestBody[key] = vals
					}
				}
			}

		}

		customWriter := &CustomResponseWriter{ResponseWriter: ctx.Writer, body: bytes.NewBufferString("")}

		ctx.Writer = customWriter

		ctx.Next()

		duration := time.Since(start)

		statusCode := ctx.Writer.Status()

		responseContentType := ctx.Writer.Header().Get("Content-Type")
		responseBodyRaw := customWriter.body.String()
		var responseBodyParse interface{}

		if strings.HasPrefix(responseContentType, "image/") {
			responseBodyParse = fmt.Sprintf("[binary data: %s]", formatSize(int64(len(responseBodyRaw))))

		} else if strings.HasPrefix(responseContentType, "application/json") ||
			strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "{") ||
			strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "[") {
			if err := json.Unmarshal([]byte(responseBodyRaw), &responseBodyParse); err != nil {
				responseBodyParse = responseBodyRaw
			}

		} else {
			responseBodyParse = responseBodyRaw
		}

		logEvent := logger.Info()
		if statusCode >= 500 {
			logEvent = logger.Error()
		} else if statusCode >= 400 {
			logEvent = logger.Warn()
		}

		logEvent.
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("query", ctx.Request.URL.RawQuery).
			Str("client_ip", ctx.ClientIP()).
			Str("user_agent", ctx.Request.UserAgent()).
			Str("referer", ctx.Request.Referer()).
			Str("protocol", ctx.Request.Proto).
			Str("host", ctx.Request.Host).
			Str("remote_addr", ctx.Request.RemoteAddr).
			Str("request_uri", ctx.Request.RequestURI).
			Int64("content_length", ctx.Request.ContentLength).
			Interface("headers", ctx.Request.Header).
			Interface("request_body", requestBody).
			Interface("response_body", responseBodyParse).
			Int("status_code", ctx.Writer.Status()).
			Int64("duration_ms", duration.Milliseconds()).
			Msg("HTTP Request Log")
	}
}

func formatSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
