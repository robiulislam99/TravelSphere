// static/js/home.js
'use strict';

(function () {
  const input = document.getElementById('home-search');
  const box   = document.getElementById('search-suggestions');
  if (!input || !box) return;

  let timer;
  let activeIndex = -1;
  let currentItems = [];

  function closeBox() {
    box.innerHTML = '';
    box.hidden = true;
    activeIndex = -1;
    currentItems = [];
  }

  function setActive(index) {
    const items = box.querySelectorAll('.suggestion-item');
    items.forEach(function (el, i) {
      el.classList.toggle('is-active', i === index);
    });
    activeIndex = index;
  }

  function renderSuggestions(countries) {
    if (!countries.length) {
      box.innerHTML =
        '<div class="suggestion-empty">No results found.</div>';
      box.hidden = false;
      return;
    }

    currentItems = countries;
    box.innerHTML = countries.map(function (c, i) {
      var capital = c.capital ? ' \u2014 ' + tsEscape(c.capital) : '';
      return '<a href="/countries/' + tsEscape(c.slug) + '" ' +
        'class="suggestion-item" ' +
        'data-index="' + i + '" ' +
        'role="option">' +
        '<span class="suggestion-name">' +
          tsEscape(c.name) +
          '<span class="suggestion-capital">' + capital + '</span>' +
        '</span>' +
      '</a>';
    }).join('');
    box.hidden = false;
  }

  function fetchSuggestions(q) {
    tsShowSpinner(box);
    box.hidden = false;

    tsAjax({ url: '/api/countries', data: { search: q } })
      .then(function (res) {
        renderSuggestions((res.data || []).slice(0, 8));
      })
      .catch(function () {
        closeBox();
      });
  }

  // Input — debounced fetch
  input.addEventListener('input', function () {
    clearTimeout(timer);
    const q = input.value.trim();
    if (q.length < 2) { closeBox(); return; }
    timer = setTimeout(function () { fetchSuggestions(q); }, 300);
  });

  // Keyboard navigation
  input.addEventListener('keydown', function (e) {
    const items = box.querySelectorAll('.suggestion-item');
    if (!items.length) return;

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActive(Math.min(activeIndex + 1, items.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActive(Math.max(activeIndex - 1, 0));
    } else if (e.key === 'Enter' && activeIndex >= 0) {
      e.preventDefault();
      items[activeIndex].click();
    } else if (e.key === 'Escape') {
      closeBox();
      input.blur();
    }
  });

  // Click outside — close
  document.addEventListener('click', function (e) {
    if (!input.contains(e.target) && !box.contains(e.target)) {
      closeBox();
    }
  });

  // Reopen on focus if value exists
  input.addEventListener('focus', function () {
    const q = input.value.trim();
    if (q.length >= 2 && !currentItems.length) {
      fetchSuggestions(q);
    }
  });
})();