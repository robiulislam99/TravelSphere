package controllers

import (
	stdc "context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

// ── mock session store ─────────────────────────────────────

type mockSessionStore struct {
	data map[interface{}]interface{}
}

func (m *mockSessionStore) Get(ctx stdc.Context, key interface{}) interface{} {
	return m.data[key]
}

func (m *mockSessionStore) Set(ctx stdc.Context, key interface{}, val interface{}) error {
	m.data[key] = val
	return nil
}

func (m *mockSessionStore) Delete(ctx stdc.Context, key interface{}) error {
	delete(m.data, key)
	return nil
}

func (m *mockSessionStore) Flush(ctx stdc.Context) error {
	m.data = make(map[interface{}]interface{})
	return nil
}

func (m *mockSessionStore) SessionID(ctx stdc.Context) string {
	return "test-session-id"
}

func (m *mockSessionStore) SessionRelease(ctx stdc.Context, w http.ResponseWriter) {
	// No-op for mock
}

func (m *mockSessionStore) SessionReleaseIfPresent(ctx stdc.Context, w http.ResponseWriter) {
}

var _ session.Store = (*mockSessionStore)(nil)

// ── helpers ────────────────────────────────────────────────

// newBaseControllerCtx wires a BaseController to a fresh request/response pair
// and returns both the controller and the response recorder for assertions.
func newBaseControllerCtx(method, path string) (*BaseController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	ctx := context.NewContext()
	ctx.Reset(w, r)

	c := &BaseController{}
	c.Ctx = ctx
	c.Data = make(map[interface{}]interface{})
	
	// Initialize session with our mock store
	c.CruSession = &mockSessionStore{
		data: make(map[interface{}]interface{}),
	}

	return c, w
}

// ── default template data ──────────────────────────────────

func TestBaseController_Prepare_DefaultTitle(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")
    c.Prepare()

    if c.Data["Title"] != "TravelSphere" {
        t.Errorf("Title = %q; want TravelSphere", c.Data["Title"])
    }
}

func TestBaseController_Prepare_DefaultMetaDescription(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")
    c.Prepare()

    meta, ok := c.Data["MetaDescription"].(string)
    if !ok || meta == "" {
        t.Error("MetaDescription should be a non-empty string")
    }
}

func TestBaseController_Prepare_DefaultActivePage(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")
    c.Prepare()

    if c.Data["ActivePage"] != "" {
        t.Errorf("ActivePage default = %q; want empty string", c.Data["ActivePage"])
    }
}

func TestBaseController_Prepare_DefaultPageScript(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")
    c.Prepare()

    if c.Data["PageScript"] != "" {
        t.Errorf("PageScript default = %q; want empty string", c.Data["PageScript"])
    }
}

// ── session: not logged in ─────────────────────────────────

func TestBaseController_Prepare_NotLoggedIn_LoggedInFalse(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")
    c.Prepare()

    loggedIn, ok := c.Data["LoggedIn"].(bool)
    if !ok {
        t.Fatal("LoggedIn should be a bool")
    }
    if loggedIn {
        t.Error("LoggedIn should be false when no session exists")
    }
}

func TestBaseController_Prepare_NotLoggedIn_FirstNameEmpty(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")
    c.Prepare()

    firstName, ok := c.Data["FirstName"].(string)
    if !ok {
        t.Fatal("FirstName should be a string")
    }
    if firstName != "" {
        t.Errorf("FirstName = %q; want empty string when not logged in", firstName)
    }
}

// ── session: logged in ─────────────────────────────────────

// setSession is a test helper that injects session values directly
// into the Beego session store on the controller's context.
func setSession(c *BaseController, key, value string) {
	c.SetSession(key, value)
}

// setSessionGeneric works with any controller type that has BaseController embedded
func setSessionGeneric(store *mockSessionStore, key, value string) {
	store.Set(stdc.Background(), key, value)
}

func TestBaseController_Prepare_LoggedIn_LoggedInTrue(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")

    // Simulate an active session
    setSession(c, "username",   "robiulislam99")
    setSession(c, "first_name", "Robiul")

    c.Prepare()

    loggedIn, ok := c.Data["LoggedIn"].(bool)
    if !ok {
        t.Fatal("LoggedIn should be a bool")
    }
    if !loggedIn {
        t.Error("LoggedIn should be true when session username is set")
    }
}

func TestBaseController_Prepare_LoggedIn_FirstNameInjected(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")

    setSession(c, "username",   "robiulislam99")
    setSession(c, "first_name", "Robiul")

    c.Prepare()

    firstName, ok := c.Data["FirstName"].(string)
    if !ok {
        t.Fatal("FirstName should be a string")
    }
    if firstName != "Robiul" {
        t.Errorf("FirstName = %q; want Robiul", firstName)
    }
}

func TestBaseController_Prepare_UsernameSetButEmpty_NotLoggedIn(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")

    // username key exists but value is empty string
    setSession(c, "username", "")

    c.Prepare()

    loggedIn, _ := c.Data["LoggedIn"].(bool)
    if loggedIn {
        t.Error("LoggedIn should be false when username session value is empty string")
    }
}

// ── data keys always present ───────────────────────────────

func TestBaseController_Prepare_AllRequiredKeysPresent(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")
    c.Prepare()

    required := []string{
        "Title", "MetaDescription", "ActivePage", "PageScript",
        "LoggedIn", "FirstName",
    }
    for _, key := range required {
        if _, exists := c.Data[key]; !exists {
            t.Errorf("Prepare() did not set required template key: %q", key)
        }
    }
}

// ── Prepare is idempotent ──────────────────────────────────

func TestBaseController_Prepare_CalledTwice_StableDefaults(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")

    c.Prepare()
    c.Prepare() // second call should not panic or corrupt data

    if c.Data["Title"] != "TravelSphere" {
        t.Errorf("Title after double Prepare = %q; want TravelSphere", c.Data["Title"])
    }
    if _, exists := c.Data["LoggedIn"]; !exists {
        t.Error("LoggedIn key missing after double Prepare()")
    }
}

// ── services.Init() is called ─────────────────────────────

func TestBaseController_Prepare_ServicesInitialised(t *testing.T) {
    c, _ := newBaseControllerCtx(http.MethodGet, "/")

    // If services.Init() panics or is not idempotent this test will fail
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("Prepare() panicked during services.Init(): %v", r)
        }
    }()

    c.Prepare()
    c.Prepare() // second call — Init() must be idempotent via sync.Once
}

// ── web.Controller embedding ───────────────────────────────

func TestBaseController_EmbedsWebController(t *testing.T) {
	c := &BaseController{}
	// Verify the embedded Controller field exists and is the right type
	var _ web.Controller = c.Controller
}