// controllers/api/attraction.go
// AttractionAPIController handles GET /api/attractions.
// Full implementation in Phase 6.
package api

import "github.com/beego/beego/v2/server/web"

// AttractionAPIController serves JSON attraction data.
type AttractionAPIController struct {
	web.Controller
}

// GetByCoords handles GET /api/attractions?lat=...&lon=...
func (c *AttractionAPIController) GetByCoords() {
	c.Data["json"] = map[string]interface{}{"data": []interface{}{}, "message": "stub — Phase 6"}
	c.ServeJSON()
}