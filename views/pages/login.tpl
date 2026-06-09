<div style="min-height:100vh; background:linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
            display:flex; align-items:center; justify-content:center; padding:1rem;">
    <div style="width:100%; max-width:420px; background:#fff;
                border-radius:16px; padding:2.5rem;
                box-shadow:0 8px 40px rgba(0,0,0,.18);">

        <!-- Logo -->
        <div style="text-align:center; margin-bottom:2rem;">
            <a href="/" style="text-decoration:none; color:#1a73e8;
                               font-size:1.5rem; font-weight:700;">
                ✈ TravelSphere
            </a>
            <p style="color:#666; font-size:.9rem; margin-top:.25rem;">
                Your personal travel companion
            </p>
        </div>

        <h2 style="margin-bottom:.4rem; color:#1a1a2e; font-size:1.4rem;">
            Welcome back
        </h2>
        <p style="color:#666; margin-bottom:1.5rem; font-size:.9rem;">
            Enter your name to access your wishlist and dashboard.
        </p>

        {{if .Error}}
        <div style="background:#fee2e2; color:#dc2626; padding:.75rem 1rem;
                    border-radius:8px; margin-bottom:1rem; font-size:.9rem;">
            {{.Error}}
        </div>
        {{end}}

        <form method="POST" action="/login">
            <div style="margin-bottom:1.25rem;">
                <label style="display:block; margin-bottom:.5rem;
                              font-weight:500; color:#1a1a2e; font-size:.9rem;">
                    Your name
                </label>
                <input
                    type="text"
                    name="username"
                    class="form-control"
                    placeholder="e.g. Robiul Islam"
                    style="width:100%; font-size:1rem;"
                    autofocus
                >
            </div>
            <button type="submit" class="btn btn--primary"
                    style="width:100%; padding:.75rem; font-size:1rem;">
                Continue →
            </button>
        </form>

        <p style="text-align:center; margin-top:1.5rem;
                  font-size:.85rem; color:#999;">
            No account needed — just enter your name.
        </p>
    </div>
</div>