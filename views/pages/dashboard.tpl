<section class="page-hero">
    <div class="container">
        <h1 class="page-hero__title">Travel Dashboard</h1>
        <p class="page-hero__subtitle">A summary of your travel plans and visited destinations.</p>
    </div>
</section>

<div class="container">

    <!-- AJAX partial: only this section refreshes after wishlist changes -->
    <div id="dashboard-stats">
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-card__number">{{.Summary.Total}}</div>
                <div class="stat-card__label">Total Saved</div>
            </div>
            <div class="stat-card">
                <div class="stat-card__number" style="color:var(--color-warning);">{{.Summary.Planned}}</div>
                <div class="stat-card__label">Planned Trips</div>
            </div>
            <div class="stat-card">
                <div class="stat-card__number" style="color:var(--color-success);">{{.Summary.Visited}}</div>
                <div class="stat-card__label">Visited</div>
            </div>
        </div>
    </div>

    <section class="section">
        <h2 class="section-title">Saved Destinations</h2>
        {{if .WishlistEntries}}
        <div class="country-grid">
            {{range .WishlistEntries}}
            <div class="card">
                <div class="card-body">
                    <div style="display:flex; justify-content:space-between; align-items:flex-start; margin-bottom:0.5rem;">
                        <strong style="font-size:1rem;">{{.CountryName}}</strong>
                        <span class="badge badge-{{.StatusClass}}">{{.Status}}</span>
                    </div>
                    {{if .Note}}
                    <p style="font-size:0.875rem; color:var(--color-text-muted);">{{.Note}}</p>
                    {{end}}
                    <p style="font-size:0.8125rem; color:var(--color-text-light); margin-top:0.5rem;">
                        Added {{.FormattedDate}}
                    </p>
                    <div style="margin-top:0.875rem;">
                        <a href="/countries/{{.Slug}}" class="btn btn-secondary btn-sm">View Details</a>
                    </div>
                </div>
            </div>
            {{end}}
        </div>
        {{else}}
        <div class="empty-state">
            <div class="empty-state__icon">🗺</div>
            <div class="empty-state__title">No destinations yet</div>
            <div class="empty-state__text">Explore countries and start building your wishlist.</div>
            <a href="/countries" class="btn btn-primary" style="margin-top:1.25rem;">Explore Countries</a>
        </div>
        {{end}}
    </section>

</div>