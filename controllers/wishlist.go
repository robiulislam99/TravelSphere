// controllers/wishlist.go
// WishlistController handles GET /wishlist (protected route).
// Full implementation with WishlistService comes in Phase 5.
package controllers

// WishlistController handles the wishlist SSR page.
type WishlistController struct {
	BaseController
}

// Prepare sets page-level template data.
func (c *WishlistController) Prepare() {
	c.BaseController.Prepare()
	c.Data["Title"] = "My Wishlist"
	c.Data["ActivePage"] = "wishlist"
}

// Get renders the wishlist page.
// TODO (Phase 5): fetch real entries from WishlistService.
func (c *WishlistController) Get() {
	c.Data["WishlistEntries"] = []interface{}{}

	c.Layout = "layout/base.tpl"
	c.TplName = "pages/wishlist.tpl"
}