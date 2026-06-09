package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

// ── DashboardController.Get() ──────────────────────────────

func TestDashboardController_Get_SetsActivePage(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")
	c.Get()

	if c.Data["ActivePage"] != "dashboard" {
		t.Errorf("ActivePage = %q; want dashboard", c.Data["ActivePage"])
	}
}

func TestDashboardController_Get_SetsTitle(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")
	c.Get()

	if c.Data["Title"] != "Dashboard" {
		t.Errorf("Title = %q; want Dashboard", c.Data["Title"])
	}
}

func TestDashboardController_Get_SetsPageScript(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")
	c.Get()

	if c.Data["PageScript"] != "dashboard.js" {
		t.Errorf("PageScript = %q; want dashboard.js", c.Data["PageScript"])
	}
}

func TestDashboardController_Get_SetsCorrectTemplate(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")
	c.Get()

	if c.TplName != "pages/dashboard.tpl" {
		t.Errorf("TplName = %q; want pages/dashboard.tpl", c.TplName)
	}
}

func TestDashboardController_Get_SetsCorrectLayout(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")
	c.Get()

	if c.Layout != "layout/base.tpl" {
		t.Errorf("Layout = %q; want layout/base.tpl", c.Layout)
	}
}

func TestDashboardController_Get_WithoutSession_EmptyUsername(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")
	// No session set
	c.Get()

	// Should still populate Summary and WishlistEntries, but for empty username
	if _, exists := c.Data["Summary"]; !exists {
		t.Error("Summary key missing from template data")
	}

	if _, exists := c.Data["WishlistEntries"]; !exists {
		t.Error("WishlistEntries key missing from template data")
	}
}

func TestDashboardController_Get_WithSession_PopulateSummaryAndWishlist(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")

	// Set session with a username
	setSession(&c.BaseController, "username", "testuser")

	c.Get()

	// Summary should be populated
	if _, exists := c.Data["Summary"]; !exists {
		t.Error("Summary key missing from template data")
	}

	// WishlistEntries should be populated
	if _, exists := c.Data["WishlistEntries"]; !exists {
		t.Error("WishlistEntries key missing from template data")
	}
}

func TestDashboardController_Prepare_CallsBaseControllerPrepare(t *testing.T) {
	c, _ := newDashboardControllerCtx(http.MethodGet, "/dashboard")
	c.Prepare()

	// BaseController.Prepare() sets these
	if c.Data["Title"] != "Dashboard" {
		t.Errorf("Title after Prepare = %q; want Dashboard", c.Data["Title"])
	}

	if c.Data["ActivePage"] != "dashboard" {
		t.Errorf("ActivePage after Prepare = %q; want dashboard", c.Data["ActivePage"])
	}
}

// ── Test helper ───────────────────────────────────────────

func newDashboardControllerCtx(method, path string) (*DashboardController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	ctx := context.NewContext()
	ctx.Reset(w, r)

	c := &DashboardController{}
	c.Ctx = ctx
	c.Data = make(map[interface{}]interface{})
	c.CruSession = &mockSessionStore{
		data: make(map[interface{}]interface{}),
	}

	// Call Prepare to initialize template data (ActivePage, Title, etc.)
	c.Prepare()

	return c, w
}
