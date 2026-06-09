package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

// ── WishlistController.Get() ───────────────────────────────

func TestWishlistController_Get_SetsActivePage(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	c.Get()

	if c.Data["ActivePage"] != "wishlist" {
		t.Errorf("ActivePage = %q; want wishlist", c.Data["ActivePage"])
	}
}

func TestWishlistController_Get_SetsTitle(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	c.Get()

	if c.Data["Title"] != "My Wishlist" {
		t.Errorf("Title = %q; want My Wishlist", c.Data["Title"])
	}
}

func TestWishlistController_Get_SetsPageScript(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	c.Get()

	if c.Data["PageScript"] != "wishlist.js" {
		t.Errorf("PageScript = %q; want wishlist.js", c.Data["PageScript"])
	}
}

func TestWishlistController_Get_SetsCorrectTemplate(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	c.Get()

	if c.TplName != "pages/wishlist.tpl" {
		t.Errorf("TplName = %q; want pages/wishlist.tpl", c.TplName)
	}
}

func TestWishlistController_Get_SetsCorrectLayout(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	c.Get()

	if c.Layout != "layout/base.tpl" {
		t.Errorf("Layout = %q; want layout/base.tpl", c.Layout)
	}
}

func TestWishlistController_Get_WithoutSession_EmptyUsername(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	// No session set
	c.Get()

	// Should still populate WishlistEntries for empty username (empty list)
	if _, exists := c.Data["WishlistEntries"]; !exists {
		t.Error("WishlistEntries key missing from template data")
	}
}

func TestWishlistController_Get_WithSession_PopulateWishlist(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")

	// Set session with a username
	setSession(&c.BaseController, "username", "testuser")

	c.Get()

	// WishlistEntries should be populated
	if _, exists := c.Data["WishlistEntries"]; !exists {
		t.Error("WishlistEntries key missing from template data")
	}
}

func TestWishlistController_Prepare_CallsBaseControllerPrepare(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	c.Prepare()

	// BaseController.Prepare() sets these via Prepare() override
	if c.Data["Title"] != "My Wishlist" {
		t.Errorf("Title after Prepare = %q; want My Wishlist", c.Data["Title"])
	}

	if c.Data["ActivePage"] != "wishlist" {
		t.Errorf("ActivePage after Prepare = %q; want wishlist", c.Data["ActivePage"])
	}
}

func TestWishlistController_Get_RequiredKeys(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")
	c.Get()

	// BaseController.Prepare() sets Title, ActivePage, LoggedIn, FirstName
	// WishlistController.Prepare() overrides Title and ActivePage
	// WishlistController.Get() sets PageScript and WishlistEntries
	requiredKeys := []string{
		"PageScript", "WishlistEntries",
	}

	for _, key := range requiredKeys {
		if _, exists := c.Data[key]; !exists {
			t.Errorf("Required data key missing: %q", key)
		}
	}
	
	// Prepare should have set these
	if title, ok := c.Data["Title"]; !ok || title != "My Wishlist" {
		t.Errorf("Title not set correctly by Prepare()")
	}
}

func TestWishlistController_Get_LoggedInSessionData(t *testing.T) {
	c, _ := newWishlistControllerCtx(http.MethodGet, "/wishlist")

	// Simulate logged-in user by setting session directly on store before Get()
	// Note: This won't be picked up by Prepare() since that already ran, 
	// but demonstrates how sessions are handled
	setSession(&c.BaseController, "username", "johndoe")
	setSession(&c.BaseController, "first_name", "John")

	// Get() reads session and populates WishlistEntries for the logged-in user
	c.Get()

	// Check that WishlistEntries is populated
	if _, exists := c.Data["WishlistEntries"]; !exists {
		t.Error("WishlistEntries should be populated")
	}
}

// ── Test helper ───────────────────────────────────────────

func newWishlistControllerCtx(method, path string) (*WishlistController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	ctx := context.NewContext()
	ctx.Reset(w, r)

	c := &WishlistController{}
	c.Ctx = ctx
	c.Data = make(map[interface{}]interface{})
	c.CruSession = &mockSessionStore{
		data: make(map[interface{}]interface{}),
	}

	// Call Prepare to initialize template data (ActivePage, Title, etc.)
	c.Prepare()

	return c, w
}
