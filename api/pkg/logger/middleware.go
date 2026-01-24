package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/fleimkeipa/lifery/model"

	"github.com/labstack/echo/v4"
)

// ResponseWriter wraps echo.Response to capture the response body
type responseWriter struct {
	body *bytes.Buffer
	echo.Response
}

// Write captures the response body while continuing to write to the original response
func (rc *responseWriter) Write(b []byte) (int, error) {
	rc.body.Write(b)
	return rc.Response.Write(b)
}

// Middleware returns a middleware that logs HTTP requests
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip logging for swagger
			if c.Path() == "/swagger/*" {
				return next(c)
			}

			start := time.Now()

			// Wrap the response writer to capture the body if needed
			// We only really care about the body for error responses (>= 400)
			// Small overhead for success responses, but allows us to log error details
			res := c.Response()
			bodyBuffer := new(bytes.Buffer)
			originalWriter := res.Writer
			writer := &responseWriter{
				Response: *res,
				body:     bodyBuffer,
			}
			res.Writer = writer

			// Execute the next handler
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// Restore the original writer (optional but good practice)
			res.Writer = originalWriter

			latency := time.Since(start)
			req := c.Request()
			status := res.Status

			// Basic fields for every log
			fields := []any{
				"method", req.Method,
				"uri", req.RequestURI,
				"status", status,
				"latency", latency.String(),
			}

			devFields := []any{
				"request_id", res.Header().Get(echo.HeaderXRequestID),
				"remote_ip", c.RealIP(),
			}

			if os.Getenv("STAGE") != model.StageProd {
				devFields = append(devFields, fields...)
				fields = devFields
			}

			// If it's an error, try to extract the error message from the captured body
			if status >= 400 {
				var errorData struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				if json.Unmarshal(bodyBuffer.Bytes(), &errorData) == nil {
					if errorData.Error != "" {
						fields = append(fields, "error_detail", errorData.Error)
					}
					if errorData.Message != "" {
						fields = append(fields, "error_message", errorData.Message)
					}
				}
			}

			// Log based on status code
			switch {
			case status >= 500:
				Log.Errorw("Server Error", fields...)
			case status >= 400:
				Log.Errorw("Client Error", fields...)
			case status >= 300:
				Log.Infow("Redirection", fields...)
			default:
				Log.Infow("Success", fields...)
			}

			return nil
		}
	}
}
