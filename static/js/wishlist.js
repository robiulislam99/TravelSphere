// static/js/wishlist.js
// AJAX edit, status change, and delete on /wishlist.
// Only #wishlist-rows changes — no full page reload.
'use strict';

(function () {
  const AUTH_TOKEN = 'travelsphere-dev-token';
  const rowsBox    = document.getElementById('wishlist-rows');
  const overlay    = document.getElementById('edit-modal-overlay');
  if (!rowsBox) return;

  // ── Delete ─────────────────────────────────────────────────
  rowsBox.addEventListener('click', function (e) {
    const delBtn = e.target.closest('.btn-delete');
    if (!delBtn) return;

    const id      = delBtn.dataset.id;
    const country = delBtn.dataset.country;
    if (!confirm('Remove ' + country + ' from your wishlist?')) return;

    delBtn.disabled = true;

    fetch('/api/wishlist/' + id, {
      method: 'DELETE',
      headers: { 'Authorization': 'Bearer ' + AUTH_TOKEN },
    })
      .then(function (res) { return res.json(); })
      .then(function () { refreshRows(); refreshStats(); })
      .catch(function () {
        alert('Delete failed. Please try again.');
        delBtn.disabled = false;
      });
  });

  // ── Open edit modal ────────────────────────────────────────
  rowsBox.addEventListener('click', function (e) {
    const editBtn = e.target.closest('.btn-edit-note');
    if (!editBtn || !overlay) return;

    document.getElementById('edit-entry-id').value   = editBtn.dataset.id;
    document.getElementById('edit-note').value        = editBtn.dataset.note || '';
    document.getElementById('edit-status').value      = editBtn.dataset.status || 'Planned';
    document.getElementById('edit-feedback').innerHTML = '';
    overlay.style.display = 'flex';
  });

  // ── Cancel edit ────────────────────────────────────────────
  const cancelBtn = document.getElementById('btn-cancel-edit');
  if (cancelBtn) cancelBtn.addEventListener('click', function () {
    overlay.style.display = 'none';
  });

  // Close overlay on backdrop click
  if (overlay) overlay.addEventListener('click', function (e) {
    if (e.target === overlay) overlay.style.display = 'none';
  });

  // ── Save edit ──────────────────────────────────────────────
  const saveBtn = document.getElementById('btn-save-edit');
  if (saveBtn) saveBtn.addEventListener('click', function () {
    const id     = document.getElementById('edit-entry-id').value;
    const note   = document.getElementById('edit-note').value;
    const status = document.getElementById('edit-status').value;
    const fb     = document.getElementById('edit-feedback');

    saveBtn.disabled = true;

    fetch('/api/wishlist/' + id, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + AUTH_TOKEN,
      },
      body: JSON.stringify({ note: note, status: status }),
    })
      .then(function (res) { return res.json().then(function (j) { return { ok: res.ok, body: j }; }); })
      .then(function (r) {
        if (!r.ok) {
          tsShowAlert(fb, r.body.message || 'Update failed.', 'error');
          saveBtn.disabled = false;
          return;
        }
        overlay.style.display = 'none';
        saveBtn.disabled = false;
        refreshRows();
        refreshStats();
      })
      .catch(function () {
        tsShowAlert(fb, 'Network error.', 'error');
        saveBtn.disabled = false;
      });
  });

  // ── Helpers ────────────────────────────────────────────────

  function badgeClass(status) {
    return status === 'Visited' ? 'badge-visited' : 'badge-planned';
  }

  function refreshRows() {
    fetch('/api/wishlist', {
      headers: { 'Authorization': 'Bearer ' + AUTH_TOKEN },
    })
      .then(function (r) { return r.json(); })
      .then(function (res) {
        const entries = res.data || [];
        if (!entries.length) {
          rowsBox.innerHTML =
            '<div class="empty-state">' +
              '<div class="empty-state__icon">🧳</div>' +
              '<div class="empty-state__title">Your wishlist is empty</div>' +
              '<a href="/countries" class="btn btn-primary" style="margin-top:1rem;">Explore Countries</a>' +
            '</div>';
          return;
        }
        var rows = entries.map(function (e) {
          return '<tr data-id="' + tsEscape(e.id) + '">' +
            '<td style="font-weight:600;">' + tsEscape(e.country_name) + '</td>' +
            '<td><span class="badge ' + badgeClass(e.status) + '">' + tsEscape(e.status) + '</span></td>' +
            '<td style="color:var(--color-text-muted);">' + (e.note ? tsEscape(e.note) : '<em>No notes</em>') + '</td>' +
            '<td style="color:var(--color-text-muted);font-size:.875rem;">' + tsEscape(e.created_at ? e.created_at.slice(0, 10) : '') + '</td>' +
            '<td>' +
              '<div style="display:flex;gap:.5rem;">' +
                '<button class="btn btn-secondary btn-sm btn-edit-note" ' +
                  'data-id="' + tsEscape(e.id) + '" ' +
                  'data-note="' + tsEscape(e.note || '') + '" ' +
                  'data-status="' + tsEscape(e.status) + '">Edit</button>' +
                '<button class="btn btn-danger btn-sm btn-delete" ' +
                  'data-id="' + tsEscape(e.id) + '" ' +
                  'data-country="' + tsEscape(e.country_name) + '">Delete</button>' +
              '</div>' +
            '</td>' +
          '</tr>';
        });

        rowsBox.innerHTML =
          '<div class="table-wrapper"><table>' +
            '<thead><tr>' +
              '<th>Destination</th><th>Status</th><th>Notes</th><th>Added</th><th>Actions</th>' +
            '</tr></thead>' +
            '<tbody>' + rows.join('') + '</tbody>' +
          '</table></div>';
      });
  }

  function refreshStats() {
    const statsBox = document.getElementById('dashboard-stats');
    if (!statsBox) return; // not on dashboard page
    fetch('/api/dashboard/summary', {
      headers: { 'Authorization': 'Bearer ' + AUTH_TOKEN },
    })
      .then(function (r) { return r.json(); })
      .then(function (res) {
        const s = res.data || {};
        statsBox.querySelector && updateStatNumbers(statsBox, s);
      });
  }

  function updateStatNumbers(box, s) {
    var numbers = box.querySelectorAll('.stat-card__number');
    if (numbers[0]) numbers[0].textContent = s.total   || 0;
    if (numbers[1]) numbers[1].textContent = s.planned || 0;
    if (numbers[2]) numbers[2].textContent = s.visited || 0;
  }
})();