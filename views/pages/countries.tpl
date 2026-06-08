<section class="page-hero">
    <div class="container">
        <h1 class="page-hero__title">Explore Countries</h1>
        <p class="page-hero__subtitle">Browse all countries or search by name and region.</p>
    </div>
</section>

<div class="container">

    <section class="section" style="margin-bottom:1.5rem;">
        <div class="search-row">
            <input
                type="text"
                id="country-search"
                class="form-control"
                placeholder="Search countries..."
                value="{{.SearchQuery}}"
                autocomplete="off"
                aria-label="Search countries"
            >
            <select id="region-filter" class="form-control" style="max-width:200px;" aria-label="Filter by region">
                <option value="">All Regions</option>
                <option value="Africa"   {{if eq .RegionFilter "Africa"}}selected{{end}}>Africa</option>
                <option value="Americas" {{if eq .RegionFilter "Americas"}}selected{{end}}>Americas</option>
                <option value="Asia"     {{if eq .RegionFilter "Asia"}}selected{{end}}>Asia</option>
                <option value="Europe"   {{if eq .RegionFilter "Europe"}}selected{{end}}>Europe</option>
                <option value="Oceania"  {{if eq .RegionFilter "Oceania"}}selected{{end}}>Oceania</option>
            </select>
        </div>
    </section>

    <!-- AJAX partial: only this div changes on search/filter -->
    <div id="country-results">
        <div class="country-grid">
            {{range .Countries}}
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
                        <span>👥 {{.FormattedPopulation}}</span>
                        <span>💱 {{.CurrencyDisplay}}</span>
                        <span>🗣 {{.LanguageDisplay}}</span>
                    </div>
                </div>
            </a>
            {{else}}
            <div class="empty-state">
                <div class="empty-state__icon">🔍</div>
                <div class="empty-state__title">No countries found</div>
                <div class="empty-state__text">Try a different search term or region.</div>
            </div>
            {{end}}
        </div>
    </div>

</div>