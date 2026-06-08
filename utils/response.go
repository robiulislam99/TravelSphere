// utils/response.go
// Standard JSON response helpers used by all API controllers.
// Ensures every error and success response has a consistent shape.
package utils

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

// APIResponse is the standard envelope for all JSON API responses.
//
// Success:  { "data": <payload>, "message": "ok" }
// Error:    { "data": null,      "message": "<reason>", "status": 400 }
type APIResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

// SendSuccess writes a 200 JSON response with the given data payload.
func SendSuccess(c *web.Controller, data interface{}) {
	c.Data["json"] = APIResponse{
		Data:    data,
		Message: "ok",
		Status:  http.StatusOK,
	}
	c.ServeJSON()
}

// SendCreated writes a 201 JSON response (used after POST /api/wishlist).
func SendCreated(c *web.Controller, data interface{}) {
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = APIResponse{
		Data:    data,
		Message: "created",
		Status:  http.StatusCreated,
	}
	c.ServeJSON()
}

// SendBadRequest writes a 400 JSON response with a validation error message.
func SendBadRequest(c *web.Controller, message string) {
	c.Ctx.Output.SetStatus(http.StatusBadRequest)
	c.Data["json"] = APIResponse{
		Data:    nil,
		Message: message,
		Status:  http.StatusBadRequest,
	}
	c.ServeJSON()
}

// SendNotFound writes a 404 JSON response.
func SendNotFound(c *web.Controller, message string) {
	c.Ctx.Output.SetStatus(http.StatusNotFound)
	c.Data["json"] = APIResponse{
		Data:    nil,
		Message: message,
		Status:  http.StatusNotFound,
	}
	c.ServeJSON()
}

// SendUnauthorized writes a 401 JSON response.
func SendUnauthorized(c *web.Controller) {
	c.Ctx.Output.SetStatus(http.StatusUnauthorized)
	c.Data["json"] = APIResponse{
		Data:    nil,
		Message: "Unauthorized: valid Bearer token required",
		Status:  http.StatusUnauthorized,
	}
	c.ServeJSON()
}

// SendInternalError writes a 500 JSON response.
// The internal error detail is NOT exposed to the client.
func SendInternalError(c *web.Controller) {
	c.Ctx.Output.SetStatus(http.StatusInternalServerError)
	c.Data["json"] = APIResponse{
		Data:    nil,
		Message: "An internal error occurred. Please try again later.",
		Status:  http.StatusInternalServerError,
	}
	c.ServeJSON()
}