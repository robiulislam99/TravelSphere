// static/js/home.js
// AJAX destination search on the home page.
// Updates #search-suggestions only — no full page reload.
'use strict';

(function () {
  const input = document.getElementById('home-search');
  const box   = document.getElementById('search-suggestions');
  if (!input || !box) return;

  let timer;

  input.addEventListener('input', function () {
    clearTimeout(timer);
    const q = input.value.trim();

    if (q.length < 2) { box.innerHTML = ''; return; }

    timer = setTimeout(function () {
      tsShowSpinner(box);

      tsAjax({ url: '/api/countries', data: { search: q } })
        .then(function (res) {
          const countries = (res.data || []).slice(0, 6);
          if (!countries.length) {
            box.innerHTML = '<p style="color:var(--color-text-muted);font-size:.9rem;">No results found.</p>';
            return;
          }
          box.innerHTML = countries.map(function (c) {
            return '<a href="/countries/' + tsEscape(c.slug) + '" ' +
              'style="display:inline-flex;align-items:center;gap:.5rem;' +
              'background:rgba(255,255,255,.15);border-radius:6px;' +
              'padding:.375rem .75rem;margin:.25rem .25rem 0 0;' +
              'color:#fff;font-size:.875rem;text-decoration:none;">' +
              tsEscape(c.flag_emoji) + ' ' + tsEscape(c.name) +
              '</a>';
          }).join('');
        })
        .catch(function () {
          box.innerHTML = '';
        });
    }, 300);
  });
})();