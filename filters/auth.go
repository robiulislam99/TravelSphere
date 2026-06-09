// filters/auth.go
// AuthFilter protects routes that require authentication.
// SSR routes (/wishlist, /dashboard) — checked via session cookie.
// API routes (/api/wishlist, /api/dashboard) — checked via Bearer token.
package filters

import (
    "net/http"
    "os"
    "strings"

    "github.com/beego/beego/v2/server/web/context"
)

func AuthFilter(ctx *context.Context) {
    path := ctx.Request.URL.Path

    if strings.HasPrefix(path, "/api/") {
        // API routes — Bearer token check
        expectedToken := os.Getenv("AUTH_TOKEN")
        if expectedToken == "" {
            expectedToken = "travelsphere-dev-token"
        }

        authHeader := ctx.Request.Header.Get("Authorization")
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == "" {
            token = ctx.Input.Query("token")
        }

        if token != expectedToken {
            ctx.Output.SetStatus(http.StatusUnauthorized)
            _ = ctx.Output.JSON(map[string]interface{}{
                "message": "Unauthorized: valid Bearer token required",
                "status":  http.StatusUnauthorized,
            }, false, false)
        }

    } else {
        // SSR routes — session check
        username := ctx.Input.Session("username")
        if username == nil || username == "" {
            ctx.Redirect(http.StatusFound, "/login")
        }
    }
}