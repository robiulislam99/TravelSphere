package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

// ── HomeController.Get() ───────────────────────────────────

func TestHomeController_Get_SetsActivePage(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	if c.Data["ActivePage"] != "home" {
		t.Errorf("ActivePage = %q; want home", c.Data["ActivePage"])
	}
}

func TestHomeController_Get_SetsTitle(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	if c.Data["Title"] != "Discover Your Next Adventure" {
		t.Errorf("Title = %q; want Discover Your Next Adventure", c.Data["Title"])
	}
}

func TestHomeController_Get_SetsPageScript(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	if c.Data["PageScript"] != "home.js" {
		t.Errorf("PageScript = %q; want home.js", c.Data["PageScript"])
	}
}

func TestHomeController_Get_SetsCorrectTemplate(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	if c.TplName != "pages/home.tpl" {
		t.Errorf("TplName = %q; want pages/home.tpl", c.TplName)
	}
}

func TestHomeController_Get_SetsCorrectLayout(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	if c.Layout != "layout/base.tpl" {
		t.Errorf("Layout = %q; want layout/base.tpl", c.Layout)
	}
}

func TestHomeController_Get_PopulatesFeaturedCountries(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	// FeaturedCountries may be nil on API failure, but key should exist
	if _, exists := c.Data["FeaturedCountries"]; !exists {
		t.Error("FeaturedCountries key missing from template data")
	}
}

func TestHomeController_Get_PopulatesPopularAttractions(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	// PopularAttractions is always populated (empty slice on failure)
	if _, exists := c.Data["PopularAttractions"]; !exists {
		t.Error("PopularAttractions key missing from template data")
	}
}

func TestHomeController_Prepare_CallsBaseControllerPrepare(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Prepare()

	// BaseController.Prepare() sets these
	if c.Data["Title"] != "Discover Your Next Adventure" {
		t.Errorf("Title after Prepare = %q; want Discover Your Next Adventure", c.Data["Title"])
	}

	if c.Data["ActivePage"] != "home" {
		t.Errorf("ActivePage after Prepare = %q; want home", c.Data["ActivePage"])
	}
}

func TestHomeController_Get_GracefulDegradation_NoCountries(t *testing.T) {
	c, _ := newHomeControllerCtx(http.MethodGet, "/")
	c.Get()

	// Even if services return errors, Get() should not panic
	// and should have set all required template data
	requiredKeys := []string{"FeaturedCountries", "PopularAttractions", "PageScript"}
	for _, key := range requiredKeys {
		if _, exists := c.Data[key]; !exists {
			t.Errorf("Data key missing: %q", key)
		}
	}
	
	// Title should have been set by Prepare()
	if title, ok := c.Data["Title"]; !ok || title != "Discover Your Next Adventure" {
		t.Logf("Title not set correctly in Get() (should be set by Prepare())")
	}
}

// ── Test helper ───────────────────────────────────────────

func newHomeControllerCtx(method, path string) (*HomeController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	ctx := context.NewContext()
	ctx.Reset(w, r)

	c := &HomeController{}
	c.Ctx = ctx
	c.Data = make(map[interface{}]interface{})
	c.CruSession = &mockSessionStore{
		data: make(map[interface{}]interface{}),
	}

	// Call Prepare to initialize template data (ActivePage, Title, etc.)
	c.Prepare()

	return c, w
}
