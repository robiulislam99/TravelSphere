<section class="page-hero">
    <div class="container">
        <h1 class="page-hero__title">My Travel Wishlist</h1>
        <p class="page-hero__subtitle">Manage your saved destinations and travel plans.</p>
    </div>
</section>

<div class="container">

    <!-- AJAX partial: only this section changes on add/edit/delete -->
    <div id="wishlist-rows">
        {{if .WishlistEntries}}
        <div class="table-wrapper">
            <table>
                <thead>
                    <tr>
                        <th>Destination</th>
                        <th>Status</th>
                        <th>Notes</th>
                        <th>Added</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .WishlistEntries}}
                    <tr data-id="{{.ID}}">
                        <td style="font-weight:600;">{{.CountryName}}</td>
                        <td>
                            <span class="badge badge-{{.StatusClass}}">{{.Status}}</span>
                        </td>
                        <td style="max-width:220px; color:var(--color-text-muted);">
                            {{if .Note}}{{.Note}}{{else}}<em>No notes</em>{{end}}
                        </td>
                        <td style="color:var(--color-text-muted); font-size:0.875rem;">{{.FormattedDate}}</td>
                        <td>
                            <div style="display:flex; gap:0.5rem; flex-wrap:wrap;">
                                <button
                                    class="btn btn-secondary btn-sm btn-edit-note"
                                    data-id="{{.ID}}"
                                    data-note="{{.Note}}"
                                    data-status="{{.Status}}"
                                >Edit</button>
                                <button
                                    class="btn btn-danger btn-sm btn-delete"
                                    data-id="{{.ID}}"
                                    data-country="{{.CountryName}}"
                                >Delete</button>
                            </div>
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{else}}
        <div class="empty-state">
            <div class="empty-state__icon">🧳</div>
            <div class="empty-state__title">Your wishlist is empty</div>
            <div class="empty-state__text">Start exploring countries and add destinations you'd like to visit.</div>
            <a href="/countries" class="btn btn-primary" style="margin-top:1.25rem;">Explore Countries</a>
        </div>
        {{end}}
    </div>

</div>

<!-- Edit modal -->
<div id="edit-modal-overlay" style="display:none; position:fixed; inset:0; background:rgba(0,0,0,0.4); z-index:200; align-items:center; justify-content:center;">
    <div class="card" style="width:100%; max-width:440px; margin:auto;">
        <div class="card-body">
            <h3 style="font-size:1.0625rem; font-weight:600; margin-bottom:1rem;">Edit Wishlist Entry</h3>
            <input type="hidden" id="edit-entry-id">

            <label style="display:block; font-size:0.875rem; font-weight:500; margin-bottom:0.375rem;">Status</label>
            <select id="edit-status" class="form-control" style="margin-bottom:1rem;">
                <option value="Planned">Planned</option>
                <option value="Visited">Visited</option>
            </select>

            <label style="display:block; font-size:0.875rem; font-weight:500; margin-bottom:0.375rem;">Note</label>
            <textarea
                id="edit-note"
                class="form-control"
                rows="3"
                placeholder="Add a personal note..."
                style="resize:vertical; margin-bottom:1rem;"
            ></textarea>

            <div style="display:flex; gap:0.75rem;">
                <button id="btn-save-edit" class="btn btn-primary" style="flex:1;">Save Changes</button>
                <button id="btn-cancel-edit" class="btn btn-secondary">Cancel</button>
            </div>
            <div id="edit-feedback" style="margin-top:0.75rem;"></div>
        </div>
    </div>
</div>