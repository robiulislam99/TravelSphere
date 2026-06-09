{{define "partials/header.tpl"}}
<header class="site-header">
    <div class="container">
        <div class="header-inner">
            <a href="/" class="site-logo">
                <span class="logo-icon">✈</span>
                <span class="logo-text">TravelSphere</span>
            </a>
            <nav class="site-nav" aria-label="Main navigation">
                <ul class="nav-list">
                    <li class="nav-item {{if eq .ActivePage "home"}}nav-item--active{{end}}">
                        <a href="/" class="nav-link">Home</a>
                    </li>
                    <li class="nav-item {{if eq .ActivePage "countries"}}nav-item--active{{end}}">
                        <a href="/countries" class="nav-link">Explore</a>
                    </li>
                    {{if .LoggedIn}}
                    <li class="nav-item {{if eq .ActivePage "wishlist"}}nav-item--active{{end}}">
                        <a href="/wishlist" class="nav-link">Wishlist</a>
                    </li>
                    <li class="nav-item {{if eq .ActivePage "dashboard"}}nav-item--active{{end}}">
                        <a href="/dashboard" class="nav-link">Dashboard</a>
                    </li>
                    {{end}}
                </ul>
            </nav>

            <!-- Auth section -->
            <div class="nav-auth">
                {{if .LoggedIn}}
                <span class="nav-username">Hi, {{.FirstName}}</span>
                <a href="/logout" class="btn btn--outline btn--sm">Logout</a>
                {{else}}
                <a href="/login" class="btn btn--primary btn--sm">Login</a>
                {{end}}
            </div>

            <button class="nav-toggle" aria-label="Toggle navigation" aria-expanded="false">
                <span class="hamburger-bar"></span>
                <span class="hamburger-bar"></span>
                <span class="hamburger-bar"></span>
            </button>
        </div>
    </div>
</header>
{{end}}