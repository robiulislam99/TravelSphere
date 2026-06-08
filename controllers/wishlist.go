// controllers/wishlist.go
// WishlistController handles GET /wishlist (protected SSR page).
package controllers

import "github.com/robiulislam99/TravelSphere/services"

type WishlistController struct {
	BaseController
}

func (c *WishlistController) Prepare() {
	c.BaseController.Prepare()
	c.Data["Title"] = "My Wishlist"
	c.Data["ActivePage"] = "wishlist"
}

func (c *WishlistController) Get() {
	entries := services.Wishlist().GetAll()
	c.Data["WishlistEntries"] = entries
	c.Data["PageScript"] = "wishlist.js"

	c.Layout = "layout/base.tpl"
	c.TplName = "pages/wishlist.tpl"
}
