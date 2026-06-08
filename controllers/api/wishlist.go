// controllers/api/wishlist.go
// WishlistAPIController handles all /api/wishlist CRUD routes.
// Full implementation in Phase 6.
package api

import "github.com/beego/beego/v2/server/web"

// WishlistAPIController serves JSON responses for wishlist CRUD.
type WishlistAPIController struct {
	web.Controller
}

// GetAll handles GET /api/wishlist
func (c *WishlistAPIController) GetAll() {
	c.Data["json"] = map[string]interface{}{"data": []interface{}{}}
	c.ServeJSON()
}

// Create handles POST /api/wishlist
func (c *WishlistAPIController) Create() {
	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]interface{}{"message": "stub — Phase 6"}
	c.ServeJSON()
}

// Update handles PUT /api/wishlist/:id
func (c *WishlistAPIController) Update() {
	c.Data["json"] = map[string]interface{}{"message": "stub — Phase 6"}
	c.ServeJSON()
}

// Delete handles DELETE /api/wishlist/:id
func (c *WishlistAPIController) Delete() {
	c.Data["json"] = map[string]interface{}{"message": "stub — Phase 6"}
	c.ServeJSON()
}