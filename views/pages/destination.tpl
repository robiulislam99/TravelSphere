<section class="page-hero">
    <div class="container">
        <div style="display:flex; align-items:center; gap:1.25rem; flex-wrap:wrap;">
            {{if .Country.FlagURL}}
            <img
                src="{{.Country.FlagURL}}"
                alt="Flag of {{.Country.Name}}"
                style="height:64px; border-radius:4px; box-shadow:0 2px 8px rgba(0,0,0,0.2);"
            >
            {{end}}
            <div>
                <h1 class="page-hero__title">{{.Country.Name}}</h1>
                <p class="page-hero__subtitle">{{.Country.Region}} · {{.Country.Subregion}}</p>
            </div>
        </div>
    </div>
</section>

<div class="container">
    <div style="display:grid; grid-template-columns:2fr 1fr; gap:2rem; align-items:start;">

        <!-- Left: country info + attractions -->
        <div>
            <section class="section">
                <h2 class="section-title">Country Information</h2>
                <div class="card">
                    <div class="card-body">
                        <table style="width:100%; font-size:0.9375rem;">
                            <tbody>
                                <tr>
                                    <td style="padding:0.5rem 0; color:var(--color-text-muted); width:40%;">Capital</td>
                                    <td style="padding:0.5rem 0; font-weight:500;">{{.Country.Capital}}</td>
                                </tr>
                                <tr>
                                    <td style="padding:0.5rem 0; color:var(--color-text-muted);">Population</td>
                                    <td style="padding:0.5rem 0; font-weight:500;">{{.Country.FormattedPopulation}}</td>
                                </tr>
                                <tr>
                                    <td style="padding:0.5rem 0; color:var(--color-text-muted);">Currency</td>
                                    <td style="padding:0.5rem 0; font-weight:500;">{{.Country.CurrencyDisplay}}</td>
                                </tr>
                                <tr>
                                    <td style="padding:0.5rem 0; color:var(--color-text-muted);">Languages</td>
                                    <td style="padding:0.5rem 0; font-weight:500;">{{.Country.LanguageDisplay}}</td>
                                </tr>
                                <tr>
                                    <td style="padding:0.5rem 0; color:var(--color-text-muted);">Region</td>
                                    <td style="padding:0.5rem 0; font-weight:500;">{{.Country.Region}}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </section>

            <section class="section">
                <h2 class="section-title">Attractions &amp; Landmarks</h2>
                <div class="attraction-grid">
                    {{range .Attractions}}
                    <div class="attraction-card card">
                        <div class="card-body">
                            <div class="attraction-card__name">{{.Name}}</div>
                            <div class="attraction-card__kind">{{.PrimaryKind}}</div>
                            {{if .Distance}}
                            <div style="font-size:0.8125rem; color:var(--color-text-muted); margin-top:0.25rem;">
                                📍 {{.FormattedDistance}}
                            </div>
                            {{end}}
                        </div>
                    </div>
                    {{else}}
                    <div class="empty-state" style="grid-column:1/-1;">
                        <div class="empty-state__icon">🗺</div>
                        <div class="empty-state__title">No attractions found</div>
                        <div class="empty-state__text">Attraction data is currently unavailable.</div>
                    </div>
                    {{end}}
                </div>
            </section>
        </div>

        <!-- Right: wishlist + weather -->
        <div>
            <section class="section">
                <div class="card">
                    <div class="card-body">
                        <h3 style="font-size:1.0625rem; font-weight:600; margin-bottom:1rem;">🧳 Add to Wishlist</h3>

                        <!-- AJAX partial: only this area changes -->
                        <div id="wishlist-feedback"></div>

                        <button
                            id="btn-add-wishlist"
                            class="btn btn-primary"
                            style="width:100%;"
                            data-country="{{.Country.Name}}"
                        >
                            Add to Wishlist
                        </button>
                        <p style="font-size:0.8125rem; color:var(--color-text-muted); margin-top:0.625rem; text-align:center;">
                            Status defaults to <em>Planned</em>. Edit from your wishlist.
                        </p>
                    </div>
                </div>
            </section>

            {{if .Weather}}
            <section class="section">
                <div class="card">
                    <div class="card-body">
                        <h3 style="font-size:1.0625rem; font-weight:600; margin-bottom:0.75rem;">⛅ Current Weather</h3>
                        <p style="font-size:1.5rem; font-weight:700;">{{.Weather.TempC}}°C</p>
                        <p style="color:var(--color-text-muted);">{{.Weather.Condition}}</p>
                        <p style="font-size:0.8125rem; color:var(--color-text-muted); margin-top:0.25rem;">
                            💧 Humidity: {{.Weather.Humidity}}%
                        </p>
                    </div>
                </div>
            </section>
            {{end}}
        </div>

    </div>

    <div style="margin-top:1rem;">
        <a href="/countries" class="btn btn-secondary">← Back to Countries</a>
    </div>
</div>