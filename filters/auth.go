// filters/auth.go
// AuthFilter protects routes that require authentication.
// It checks for a valid Bearer token in the Authorization header.
// Protected SSR routes: /wishlist, /dashboard
// Protected API routes:  /api/wishlist, /api/dashboard/*
//
// For this assessment a static token is sufficient.
// Set the expected token via the AUTH_TOKEN environment variable.
// Default dev token: "travelsphere-dev-token"
package filters

import (
	"net/http"
	"os"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// AuthFilter validates the Authorization header for protected routes.
func AuthFilter(ctx *context.Context) {
	expectedToken := os.Getenv("AUTH_TOKEN")
	if expectedToken == "" {
		// Fall back to default dev token when env var not set
		expectedToken = "travelsphere-dev-token"
	}

	authHeader := ctx.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Also accept token via ?token= query param for browser SSR routes
	if token == "" {
		token = ctx.Input.Query("token")
	}

	if token != expectedToken {
		// Determine response format based on route
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
			// JSON 401 for API routes
			ctx.Output.SetStatus(http.StatusUnauthorized)
			_ = ctx.Output.JSON(map[string]interface{}{
				"message": "Unauthorized: valid Bearer token required",
				"status":  http.StatusUnauthorized,
			}, false, false)
		} else {
			// Redirect SSR routes to home with an error message
			ctx.Redirect(http.StatusFound, "/?error=unauthorized")
		}
	}
}