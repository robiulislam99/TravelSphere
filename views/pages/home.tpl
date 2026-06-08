<section class="page-hero">
    <div class="container">
        <h1 class="page-hero__title">Discover Your Next Adventure</h1>
        <p class="page-hero__subtitle">Explore countries, find attractions, and build your travel wishlist.</p>

        <div class="search-row" style="margin-top:1.5rem;">
            <div style="position:relative; max-width:480px;">
                <input
                    type="text"
                    id="home-search"
                    class="form-control"
                    placeholder="Search destinations..."
                    autocomplete="off"
                    aria-label="Search destinations"
                    style="font-size:1rem; width:100%;"
                >
                <div id="search-suggestions" class="search-suggestions" hidden></div>
            </div>
        </div>

    </div>
</section>

<div class="container">

    <section class="section">
        <h2 class="section-title">Featured Countries</h2>
        <p class="section-subtitle">Popular destinations to inspire your next trip.</p>
        <div class="country-grid">
            {{range .FeaturedCountries}}
            <a href="/countries/{{.Slug}}" class="country-card card">
                {{if .FlagURL}}
                <img class="country-card__flag" src="{{.FlagURL}}" alt="Flag of {{.Name}}" loading="lazy">
                {{else}}
                <div class="country-card__flag--placeholder">🌍</div>
                {{end}}
                <div class="country-card__body">
                    <div class="country-card__name">{{.Name}}</div>
                    <div class="country-card__meta">
                        <span>🏙 {{.Capital}}</span>
                        <span>🌍 {{.Region}}</span>
                    </div>
                </div>
            </a>
            {{else}}
            <div class="empty-state">
                <div class="empty-state__icon">🌐</div>
                <div class="empty-state__title">No countries loaded</div>
                <div class="empty-state__text">Could not fetch featured destinations right now.</div>
            </div>
            {{end}}
        </div>
    </section>

    <section class="section">
        <h2 class="section-title">Popular Attractions</h2>
        <p class="section-subtitle">Top-rated spots travellers love worldwide.</p>
        <div class="attraction-grid">
            {{range .PopularAttractions}}
            <div class="attraction-card card">
                <div class="card-body">
                    <div class="attraction-card__name">{{.Name}}</div>
                    <div class="attraction-card__kind">{{.Kinds}}</div>
                    {{if .CountryName}}
                    <div style="font-size:0.8125rem; color:var(--color-text-muted); margin-top:0.25rem;">
                        📍 {{.CountryName}}
                    </div>
                    {{end}}
                </div>
            </div>
            {{else}}
            <div class="empty-state">
                <div class="empty-state__icon">🗺</div>
                <div class="empty-state__title">No attractions loaded</div>
                <div class="empty-state__text">Could not fetch attractions right now.</div>
            </div>
            {{end}}
        </div>
    </section>

</div>