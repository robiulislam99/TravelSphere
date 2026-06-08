// TravelSphere — main.js
// Shared utilities loaded on every page
'use strict';

// ── Mobile nav toggle ──────────────────────────────────────
(function () {
    const toggle = document.querySelector('.nav-toggle');
    const nav    = document.querySelector('.site-nav');
    if (!toggle || !nav) return;

    toggle.addEventListener('click', function () {
        const isOpen = nav.classList.toggle('is-open');
        toggle.setAttribute('aria-expanded', String(isOpen));
    });

    document.addEventListener('click', function (e) {
        if (!nav.contains(e.target) && !toggle.contains(e.target)) {
            nav.classList.remove('is-open');
            toggle.setAttribute('aria-expanded', 'false');
        }
    });
})();

// ── Shared AJAX helper ─────────────────────────────────────
window.tsAjax = function (opts) {
    const method = (opts.method || 'GET').toUpperCase();
    let url = opts.url;

    if (method === 'GET' && opts.data) {
        const params = new URLSearchParams(opts.data);
        url = url + '?' + params.toString();
    }

    const headers = {};

    // Only set Content-Type for requests that have a body
    if (method !== 'GET') {
        headers['Content-Type'] = 'application/json';
    }

    const authMeta = document.querySelector('meta[name="auth-token"]');
    if (authMeta) {
        headers['Authorization'] = 'Bearer ' + authMeta.content;
    }

    const fetchOpts = { method: method, headers: headers };

    if (method !== 'GET' && opts.data) {
        fetchOpts.body = JSON.stringify(opts.data);
    }

    // Abort after 10 seconds to prevent infinite spinner
    const controller = new AbortController();
    const timeoutId  = setTimeout(function () { controller.abort(); }, 10000);
    fetchOpts.signal = controller.signal;

    return fetch(url, fetchOpts)
        .then(function (res) {
            clearTimeout(timeoutId);
            return res.json().then(function (json) {
                if (!res.ok) {
                    const err     = new Error(json.message || 'Request failed');
                    err.status    = res.status;
                    err.data      = json;
                    throw err;
                }
                return json;
            });
        })
        .catch(function (err) {
            clearTimeout(timeoutId);
            // Distinguish abort/timeout from other network errors
            if (err.name === 'AbortError') {
                throw new Error('Request timed out. Please try again.');
            }
            throw err;
        });
};

// ── Spinner helper ─────────────────────────────────────────
window.tsShowSpinner = function (container) {
    container.innerHTML =
        '<div class="loading-state"><span class="spinner"></span> Loading...</div>';
};

// ── Alert helper ───────────────────────────────────────────
// type: 'success' | 'error'
window.tsShowAlert = function (container, message, type) {
    type = type || 'success';
    container.innerHTML =
        '<div class="alert alert--' + type + '">' + escapeHtml(message) + '</div>';
};

// ── HTML escape utility ────────────────────────────────────
function escapeHtml(str) {
    const div = document.createElement('div');
    div.appendChild(document.createTextNode(String(str)));
    return div.innerHTML;
}
window.tsEscape = escapeHtml;