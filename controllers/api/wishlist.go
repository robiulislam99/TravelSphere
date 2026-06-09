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

// getUsername extracts the session username or returns empty string.
func (c *WishlistAPIController) getUsername() string {
    u := c.GetSession("username")
    if u == nil {
        return ""
    }
    return u.(string)
}

func (c *WishlistAPIController) GetAll() {
    entries := services.Wishlist().GetAll(c.getUsername())
    utils.SendSuccess(&c.Controller, entries)
}

func (c *WishlistAPIController) Create() {
    var req models.CreateWishlistRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
        utils.SendBadRequest(&c.Controller, "invalid JSON body")
        return
    }

    username := c.getUsername()
    isDup    := services.Wishlist().IsDuplicate(username, req.CountryName)
    entry, err := services.Wishlist().Create(username, &req)
    if err != nil {
        utils.SendBadRequest(&c.Controller, err.Error())
        return
    }

    utils.SendCreated(&c.Controller, map[string]interface{}{
        "entry":     entry,
        "duplicate": isDup,
    })
}

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

    entry, err := services.Wishlist().Update(c.getUsername(), id, &req)
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

func (c *WishlistAPIController) Delete() {
    id := c.Ctx.Input.Param(":id")
    if id == "" {
        utils.SendBadRequest(&c.Controller, "id is required")
        return
    }

    if err := services.Wishlist().Delete(c.getUsername(), id); err != nil {
        utils.SendNotFound(&c.Controller, err.Error())
        return
    }
    utils.SendSuccess(&c.Controller, map[string]string{"deleted": id})
}