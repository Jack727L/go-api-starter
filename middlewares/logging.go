package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

var logExcludedPaths = map[string]bool{
	"/healthz":     true,
	"/healthcheck": true,
	"/readyz":      true,
	"/livez":       true,
}

// LoggingMiddleware logs each request to stdout in a structured format.
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		err := c.Next()

		status := c.Response().StatusCode()
		duration := time.Since(start)

		if logExcludedPaths[path] {
			return err
		}

		level := "INFO"
		if status >= 500 {
			level = "ERROR"
		} else if status >= 400 {
			level = "WARN"
		}

		// Structured log line suitable for Grafana / Loki
		userID := ""
		if uid := c.Locals("userID"); uid != nil {
			userID = fmt.Sprintf(" user=%v", uid)
		}

		errContext := ""
		if errMsg := c.Locals("_err"); errMsg != nil {
			errFile, _ := c.Locals("_err_file").(string)
			errLine, _ := c.Locals("_err_line").(int)
			errFunc, _ := c.Locals("_err_func").(string)
			errContext = fmt.Sprintf(" err_file=%s err_line=%d err_func=%s error=%q", errFile, errLine, errFunc, errMsg)
		}

		fmt.Printf("%s %s %s %s %d %dms ip=%s%s%s\n",
			level, time.Now().Format(time.RFC3339),
			method, path, status, duration.Milliseconds(),
			ip, userID, errContext,
		)

		return err
	}
}
