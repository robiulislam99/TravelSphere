// controllers/api/wishlist.go
// WishlistAPIController — full CRUD JSON endpoints for the wishlist.
package api

import (
	"encoding/json"

	"github.com/beego/beego/v2/server/web"
	"github.com/robiulislam99/TravelSphere/models"
	"github.com/robiulislam99/TravelSphere/services"
	"github.com/robiulislam99/TravelSphere/utils"
)

type WishlistAPIController struct {
	web.Controller
}

// GetAll handles GET /api/wishlist
func (c *WishlistAPIController) GetAll() {
	entries := services.Wishlist().GetAll()
	utils.SendSuccess(&c.Controller, entries)
}

// Create handles POST /api/wishlist
func (c *WishlistAPIController) Create() {
	var req models.CreateWishlistRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendBadRequest(&c.Controller, "invalid JSON body")
		return
	}

	// Warn client about duplicates but still allow creation
	isDup := services.Wishlist().IsDuplicate(req.CountryName)

	entry, err := services.Wishlist().Create(&req)
	if err != nil {
		utils.SendBadRequest(&c.Controller, err.Error())
		return
	}

	resp := map[string]interface{}{
		"entry":     entry,
		"duplicate": isDup,
	}
	utils.SendCreated(&c.Controller, resp)
}

// Update handles PUT /api/wishlist/:id
func (c *WishlistAPIController) Update() {
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		utils.SendBadRequest(&c.Controller, "id is required")
		return
	}

	var req models.UpdateWishlistRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendBadRequest(&c.Controller, "invalid JSON body")
		return
	}

	entry, err := services.Wishlist().Update(id, &req)
	if err != nil {
		if err.Error() == "wishlist entry not found" {
			utils.SendNotFound(&c.Controller, err.Error())
			return
		}
		utils.SendBadRequest(&c.Controller, err.Error())
		return
	}
	utils.SendSuccess(&c.Controller, entry)
}

// Delete handles DELETE /api/wishlist/:id
func (c *WishlistAPIController) Delete() {
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		utils.SendBadRequest(&c.Controller, "id is required")
		return
	}

	if err := services.Wishlist().Delete(id); err != nil {
		utils.SendNotFound(&c.Controller, err.Error())
		return
	}
	utils.SendSuccess(&c.Controller, map[string]string{"deleted": id})
}