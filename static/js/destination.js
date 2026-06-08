// static/js/destination.js
// AJAX Add to Wishlist on /countries/:slug.
// Only #wishlist-feedback changes — rest of the page is untouched.
'use strict';

(function () {
  const btn      = document.getElementById('btn-add-wishlist');
  const feedback = document.getElementById('wishlist-feedback');
  if (!btn || !feedback) return;

  const AUTH_TOKEN = 'travelsphere-dev-token'; // matches AUTH_TOKEN env default

  btn.addEventListener('click', function () {
    const countryName = btn.dataset.country;
    btn.disabled = true;
    tsShowSpinner(feedback);

    fetch('/api/wishlist', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + AUTH_TOKEN,
      },
      body: JSON.stringify({ country_name: countryName, status: 'Planned' }),
    })
      .then(function (res) { return res.json().then(function (j) { return { ok: res.ok, body: j }; }); })
      .then(function (r) {
        if (!r.ok) {
          tsShowAlert(feedback, r.body.message || 'Could not add to wishlist.', 'error');
          btn.disabled = false;
          return;
        }
        const isDup = r.body.data && r.body.data.duplicate;
        const msg   = isDup
          ? countryName + ' is already in your wishlist.'
          : '✓ ' + countryName + ' added to wishlist!';
        tsShowAlert(feedback, msg, isDup ? 'error' : 'success');
        btn.textContent = isDup ? 'Already Added' : '✓ Added';
        // keep disabled so user can't double-add
      })
      .catch(function () {
        tsShowAlert(feedback, 'Network error. Please try again.', 'error');
        btn.disabled = false;
      });
  });
})();