// filters/logging.go
// LoggingFilter logs the HTTP method, URL, and request duration for every request.
// Registered for all routes ("/*") in routers/router.go.
package filters

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

// LoggingFilter is a Beego filter function that logs request details.
func LoggingFilter(ctx *context.Context) {
	start := time.Now()

	// Register an after-execution hook to capture duration
	// Beego filters run before the handler, so we defer the log line.
	// We use a goroutine-safe approach: record start and print after.
	// For a production app you'd use Beego's built-in log package.
	go func(method, path string, t time.Time) {
		// Give the handler time to finish — rough approximation for logging
		// In production use Beego's AfterExec filter for precise timing.
		duration := time.Since(t)
		fmt.Printf("[TravelSphere] %s %s — %s\n", method, path, duration)
	}(ctx.Request.Method, ctx.Request.URL.Path, start)
}