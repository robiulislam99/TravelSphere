// TravelSphere — main.js
// Shared utilities loaded on every page

'use strict';

// ── Mobile nav toggle ──────────────────────────────────────
(function () {
    const toggle = document.querySelector('.nav-toggle');
    const nav = document.querySelector('.site-nav');
    if (!toggle || !nav) return;

    toggle.addEventListener('click', function () {
        const isOpen = nav.classList.toggle('is-open');
        toggle.setAttribute('aria-expanded', String(isOpen));
    });

    // Close nav when clicking outside
    document.addEventListener('click', function (e) {
        if (!nav.contains(e.target) && !toggle.contains(e.target)) {
            nav.classList.remove('is-open');
            toggle.setAttribute('aria-expanded', 'false');
        }
    });
})();

// ── Shared AJAX helper ─────────────────────────────────────
// Usage:
//   tsAjax({ url: '/api/countries', method: 'GET', data: {search: 'bang'} })
//     .then(json => console.log(json))
//     .catch(err => console.error(err));
window.tsAjax = function (opts) {
    const method = (opts.method || 'GET').toUpperCase();
    let url = opts.url;

    // Append query params for GET requests
    if (method === 'GET' && opts.data) {
        const params = new URLSearchParams(opts.data);
        url = url + '?' + params.toString();
    }

    const fetchOpts = {
        method: method,
        headers: { 'Content-Type': 'application/json' },
    };

    // Add auth token from meta tag if present
    const authMeta = document.querySelector('meta[name="auth-token"]');
    if (authMeta) {
        fetchOpts.headers['Authorization'] = 'Bearer ' + authMeta.content;
    }

    // Attach body for non-GET requests
    if (method !== 'GET' && opts.data) {
        fetchOpts.body = JSON.stringify(opts.data);
    }

    return fetch(url, fetchOpts).then(function (res) {
        return res.json().then(function (json) {
            if (!res.ok) {
                // Attach server message to the thrown error
                const err = new Error(json.message || 'Request failed');
                err.status = res.status;
                err.data = json;
                throw err;
            }
            return json;
        });
    });
};

// ── Spinner helper ─────────────────────────────────────────
// Shows a loading spinner inside a target container
window.tsShowSpinner = function (container) {
    container.innerHTML =
        '<div class="loading-state"><span class="spinner"></span> Loading...</div>';
};

// ── Alert helper ───────────────────────────────────────────
// Injects a success or error alert into a target container
// type: 'success' | 'error'
window.tsShowAlert = function (container, message, type) {
    type = type || 'success';
    container.innerHTML =
        '<div class="alert alert-' + type + '">' + escapeHtml(message) + '</div>';
};

// ── HTML escape utility ────────────────────────────────────
function escapeHtml(str) {
    const div = document.createElement('div');
    div.appendChild(document.createTextNode(String(str)));
    return div.innerHTML;
}
window.tsEscape = escapeHtml;