// static/js/countries.js
// AJAX search + region filter on /countries.
// Only #country-results changes — header, nav, search box stay untouched.
'use strict';

(function () {
  const searchInput  = document.getElementById('country-search');
  const regionSelect = document.getElementById('region-filter');
  const resultsBox   = document.getElementById('country-results');
  if (!searchInput || !regionSelect || !resultsBox) return;

  let timer;

  function buildCardHTML(c) {
    var flag = c.flag_url
      ? '<img class="country-card__flag" src="' + tsEscape(c.flag_url) +
        '" alt="Flag of ' + tsEscape(c.name) + '" loading="lazy">'
      : '<div class="country-card__flag--placeholder">' + tsEscape(c.flag_emoji || '🌍') + '</div>';

    return '<a href="/countries/' + tsEscape(c.slug) + '" class="country-card card">' +
      flag +
      '<div class="country-card__body">' +
        '<div class="country-card__name">' + tsEscape(c.name) + '</div>' +
        '<div class="country-card__meta">' +
          '<span>🏙 ' + tsEscape(c.capital || '—') + '</span>' +
          '<span>👥 ' + tsEscape(c.formatted_population || '—') + '</span>' +
          '<span>💱 ' + tsEscape(c.currency_display || '—') + '</span>' +
          '<span>🗣 ' + tsEscape(c.language_display || '—') + '</span>' +
        '</div>' +
      '</div>' +
    '</a>';
  }

  function fetchCountries() {
    const q      = searchInput.value.trim();
    const region = regionSelect.value;

    tsShowSpinner(resultsBox);

    tsAjax({ url: '/api/countries', data: { search: q, region: region } })
      .then(function (res) {
        const list = res.data || [];
        if (!list.length) {
          resultsBox.innerHTML =
            '<div class="empty-state">' +
              '<div class="empty-state__icon">🔍</div>' +
              '<div class="empty-state__title">No countries found</div>' +
              '<div class="empty-state__text">Try a different search term or region.</div>' +
            '</div>';
          return;
        }
        resultsBox.innerHTML =
          '<div class="country-grid">' +
            list.map(buildCardHTML).join('') +
          '</div>';
      })
      .catch(function () {
        tsShowAlert(resultsBox, 'Could not load countries. Please try again.', 'error');
      });
  }

  // Debounced search input
  searchInput.addEventListener('input', function () {
    clearTimeout(timer);
    timer = setTimeout(fetchCountries, 300);
  });

  // Immediate region change
  regionSelect.addEventListener('change', fetchCountries);
})();