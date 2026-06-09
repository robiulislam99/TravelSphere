package controllers

import "github.com/robiulislam99/TravelSphere/services"

type WishlistController struct {
    BaseController
}

func (c *WishlistController) Prepare() {
    c.BaseController.Prepare()
    c.Data["Title"]      = "My Wishlist"
    c.Data["ActivePage"] = "wishlist"
}

func (c *WishlistController) Get() {
    username := ""
    if u := c.GetSession("username"); u != nil {
        username = u.(string)
    }

    entries := services.Wishlist().GetAll(username)
    c.Data["WishlistEntries"] = entries
    c.Data["PageScript"]      = "wishlist.js"
    c.Layout                  = "layout/base.tpl"
    c.TplName                 = "pages/wishlist.tpl"
}