// static/js/dashboard.js
// AJAX stats refresh on /dashboard.
// Only #dashboard-stats changes after wishlist mutations.
'use strict';

(function () {
  const AUTH_TOKEN = 'travelsphere-dev-token';
  const statsBox   = document.getElementById('dashboard-stats');
  if (!statsBox) return;

  // Auto-refresh stats every 30 seconds while dashboard is open
  setInterval(refreshStats, 30000);

  function refreshStats() {
    fetch('/api/dashboard/summary', {
      headers: { 'Authorization': 'Bearer ' + AUTH_TOKEN },
    })
      .then(function (r) { return r.json(); })
      .then(function (res) {
        const s = res.data || {};
        const numbers = statsBox.querySelectorAll('.stat-card__number');
        if (numbers[0]) numbers[0].textContent = s.total   || 0;
        if (numbers[1]) numbers[1].textContent = s.planned || 0;
        if (numbers[2]) numbers[2].textContent = s.visited || 0;
      })
      .catch(function () {
        // Silent fail — stale numbers are fine
      });
  }
})();