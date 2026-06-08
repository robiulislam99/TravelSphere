{{define "partials/footer.tpl"}}
<footer class="site-footer">
    <div class="container">
        <div class="footer-inner">

            <div class="footer-brand">
                <span class="logo-icon">✈</span>
                <span class="logo-text">TravelSphere</span>
                <p class="footer-tagline">Discover the world, one destination at a time.</p>
            </div>

            <nav class="footer-nav" aria-label="Footer navigation">
                <ul class="footer-nav-list">
                    <li><a href="/" class="footer-link">Home</a></li>
                    <li><a href="/countries" class="footer-link">Explore Countries</a></li>
                    <li><a href="/wishlist" class="footer-link">My Wishlist</a></li>
                    <li><a href="/dashboard" class="footer-link">Dashboard</a></li>
                </ul>
            </nav>

            <div class="footer-credits">
                <p class="footer-credit-text">Powered by</p>
                <ul class="footer-credit-list">
                    <li><a href="https://restcountries.com" target="_blank" rel="noopener" class="footer-link">REST Countries</a></li>
                    <li><a href="https://opentripmap.io" target="_blank" rel="noopener" class="footer-link">OpenTripMap</a></li>
                </ul>
            </div>

        </div>
        <div class="footer-bottom">
            <p class="footer-copy">&copy; 2025 TravelSphere. Built with Beego &amp; Go.</p>
        </div>
    </div>
</footer>
{{end}}