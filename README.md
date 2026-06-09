# TravelSphere

A travel discovery web app built with Go and Beego. Explore countries, attractions, weather, and manage a personal wishlist.

---

## Setup

### Prerequisites
- Go 1.21+
- Bee tool: `go install github.com/beego/bee/v2@latest`

### Clone & Run

```bash
# 1. Clone the repository
git clone https://github.com/robiulislam99/TravelSphere.git
cd TravelSphere

# 2. Install dependencies
go mod tidy

# 3. Set up environment variables
cp .env.example .env
# Edit .env and add your API keys (see Environment Variables section)

# 4. Run the app
bee run
```

App runs at: http://localhost:8080

---

## Environment Variables

Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

| Variable | Required | Notes |
|---|---|---|
| `OPENTRIPMAP_API_KEY` | ✅ | Get free key at https://opentripmap.io |
| `WEATHER_API_KEY` | ⬜ Optional | Get free key at https://www.weatherapi.com. Weather section hidden if not set. |
| `AUTH_TOKEN` | ✅ | Bearer token for API routes. Defaults to `travelsphere-dev-token` |

---

## Project Structure

```
TravelSphere/
├── conf/app.conf
├── controllers/
│   ├── api/
│   │   ├── attraction.go         # GET /api/attractions
│   │   ├── attraction_test.go
│   │   ├── country.go            # GET /api/countries, /api/countries/:slug
│   │   ├── country_test.go
│   │   ├── dashboard.go          # GET /api/dashboard/summary
│   │   ├── dashboard_test.go
│   │   ├── wishlist.go           # GET/POST/PUT/DELETE /api/wishlist
│   │   └── wishlist_test.go
│   ├── auth.go                   # GET/POST /login, GET /logout
│   ├── auth_test.go
│   ├── base.go                   # Shared Prepare() — session + defaults
│   ├── base_test.go
│   ├── country.go                # GET /countries, /countries/:slug
│   ├── country_test.go
│   ├── dashboard.go              # GET /dashboard
│   ├── dashboard_test.go
│   ├── home.go                   # GET /
│   ├── home_test.go
│   ├── wishlist.go               # GET /wishlist
│   └── wishlist_test.go
├── filters/
│   ├── auth.go                   # Session auth (SSR) + Bearer token (API)
│   └── logging.go
├── models/
│   ├── attraction.go
│   ├── attraction_test.go
│   ├── country.go
│   ├── dashboard.go
│   ├── dashboard_test.go
│   ├── user.go
│   ├── user_test.go
│   ├── wishlist.go
│   └── wishlist_test.go
├── routers/
│   ├── router.go                 # SSR routes + filters
│   └── api.go                    # /api/* routes
├── services/
│   ├── attraction_service.go
│   ├── attraction_service_test.go
│   ├── country_service.go
│   ├── country_service_test.go
│   ├── dashboard_service.go
│   ├── dashboard_service_test.go
│   ├── registry.go               # Singleton wiring via sync.Once
│   ├── registry_test.go
│   ├── weather_service.go
│   ├── weather_service_test.go
│   ├── wishlist_service.go
│   └── wishlist_service_test.go
├── utils/
│   ├── formatter.go
│   ├── opentripmap_client.go
│   ├── response.go
│   ├── rest_countries_client.go
│   └── validator.go
├── views/
│   ├── layout/
│   │   ├── auth.tpl              # Minimal layout for login page
│   │   └── base.tpl              # Main layout with header/footer
│   ├── pages/
│   │   ├── 404.tpl
│   │   ├── countries.tpl
│   │   ├── dashboard.tpl
│   │   ├── destination.tpl
│   │   ├── home.tpl
│   │   ├── login.tpl
│   │   └── wishlist.tpl
│   └── partials/
│       ├── header.tpl
│       └── footer.tpl
├── static/
│   ├── css/main.css
│   └── js/
│       ├── main.js               # tsAjax, tsShowSpinner, tsShowAlert
│       ├── home.js               # Search autocomplete
│       ├── countries.js          # Search + region filter
│       ├── destination.js        # Add to wishlist
│       ├── wishlist.js           # Edit/delete entries
│       └── dashboard.js          # 30s stats refresh
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── main.go
└── README.md
```

---

## Wishlist Storage

Uses **Option 1 — In-memory store with JSON API endpoints**.

- Stored in a thread-safe Go map: `map[username]map[id]*WishlistEntry`
- Each user has their own isolated wishlist keyed by session username
- Full CRUD via `/api/wishlist` endpoints
- No database or ORM used
- Data clears on server restart (intentional)

---

## Country Slug Format

Lowercase name, spaces replaced by hyphens.
Example: `United States` → `/countries/united-states`

---

## Running Tests

### Results

| Package | Coverage |
|---|---|
| `controllers` | 94.5% |
| `services` | 71.6% |
| `models` | 100% |

### Commands

```bash
go test -v -cover ./controllers
go test -v -cover ./models
go test -v -cover ./services
```

Or run all at once:

```bash
go test ./...
```

---

## Pages

| Route | Page | Auth |
|---|---|---|
| `/` | Home — featured countries + attractions | Public |
| `/countries` | Explore — search + region filter | Public |
| `/countries/:slug` | Destination — country detail, attractions, weather, wishlist | Public |
| `/wishlist` | Personal wishlist | Login required |
| `/dashboard` | Stats summary | Login required |
| `/login` | Name-based login — enter name, session created automatically | Public |

---

## AJAX Behaviour

| Page | Updates |
|---|---|
| Home | `#search-suggestions` — autocomplete dropdown |
| Explore | `#country-results` — country grid |
| Destination | `#wishlist-feedback` — add to wishlist button/message |
| Wishlist | `#wishlist-rows` — edit/delete entries |
| Dashboard | `#dashboard-stats` — auto-refreshes every 30s |

---

## Credits
- [REST Countries](https://restcountries.com)
- [OpenTripMap](https://opentripmap.io)
- [WeatherAPI](https://www.weatherapi.com)
