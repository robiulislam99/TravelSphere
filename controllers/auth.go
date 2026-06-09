package controllers

import (
    "strings"
)

type AuthController struct {
    BaseController
}

func (c *AuthController) ShowLogin() {
    if c.GetSession("username") != nil {
        c.Redirect("/wishlist", 302)
        return
    }
    c.Data["Title"]      = "Login"
    c.Data["ActivePage"] = "login"
    c.Data["Error"]      = c.GetString("error")
    c.Layout             = "layout/auth.tpl"
    c.TplName            = "pages/login.tpl"
}

func (c *AuthController) DoLogin() {
    username := strings.TrimSpace(c.GetString("username"))
    if username == "" {
        c.Redirect("/login?error=Please+enter+your+name", 302)
        return
    }

    parts     := strings.Fields(username)
    firstName := parts[0]

    c.SetSession("username",   username)
    c.SetSession("first_name", firstName)

    c.Redirect("/wishlist", 302)
}

func (c *AuthController) Logout() {
    c.DestroySession()
    c.Redirect("/", 302)
}