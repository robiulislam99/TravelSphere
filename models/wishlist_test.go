package models

import (
	"testing"
	"time"
)

func TestWishlistEntryStatusClass(t *testing.T) {
	tests := []struct {
		name     string
		status   WishlistStatus
		expected string
	}{
		{
			name:     "planned status returns lowercase",
			status:   StatusPlanned,
			expected: "planned",
		},
		{
			name:     "visited status returns lowercase",
			status:   StatusVisited,
			expected: "visited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WishlistEntry{
				Status: tt.status,
			}
			result := w.StatusClass()
			if result != tt.expected {
				t.Errorf("StatusClass() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestWishlistEntryFormattedDate(t *testing.T) {
	// Create a specific date for testing
	testDate := time.Date(2023, time.January, 15, 10, 30, 0, 0, time.UTC)
	w := &WishlistEntry{
		CreatedAt: testDate,
	}

	result := w.FormattedDate()
	expected := "15 Jan 2023"

	if result != expected {
		t.Errorf("FormattedDate() = %s, want %s", result, expected)
	}
}

func TestWishlistEntryFormattedDateMultipleDates(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "January date",
			date:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: "1 Jan 2023",
		},
		{
			name:     "December date",
			date:     time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
			expected: "25 Dec 2023",
		},
		{
			name:     "February date",
			date:     time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: "29 Feb 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WishlistEntry{
				CreatedAt: tt.date,
			}
			result := w.FormattedDate()
			if result != tt.expected {
				t.Errorf("FormattedDate() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestCreateWishlistRequestValidateSuccess(t *testing.T) {
	tests := []struct {
		name    string
		request *CreateWishlistRequest
		wantErr bool
	}{
		{
			name: "valid request with Planned status",
			request: &CreateWishlistRequest{
				CountryName: "France",
				Note:        "Visit Paris",
				Status:      "Planned",
			},
			wantErr: false,
		},
		{
			name: "valid request with Visited status",
			request: &CreateWishlistRequest{
				CountryName: "Japan",
				Note:        "Already visited",
				Status:      "Visited",
			},
			wantErr: false,
		},
		{
			name: "valid request without note",
			request: &CreateWishlistRequest{
				CountryName: "Germany",
				Status:      "Planned",
			},
			wantErr: false,
		},
		{
			name: "valid request without status (defaults to Planned)",
			request: &CreateWishlistRequest{
				CountryName: "Italy",
				Note:        "Test note",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Check default status is set
			if tt.request.Status == "" {
				t.Errorf("Status should be set to default 'Planned', got empty string")
			}
		})
	}
}

func TestCreateWishlistRequestValidateErrors(t *testing.T) {
	tests := []struct {
		name      string
		request   *CreateWishlistRequest
		wantErr   bool
		errSubstr string
	}{
		{
			name: "missing country name",
			request: &CreateWishlistRequest{
				CountryName: "",
				Status:      "Planned",
			},
			wantErr:   true,
			errSubstr: "country_name is required",
		},
		{
			name: "country name with only spaces",
			request: &CreateWishlistRequest{
				CountryName: "   ",
				Status:      "Planned",
			},
			wantErr:   true,
			errSubstr: "country_name is required",
		},
		{
			name: "invalid status",
			request: &CreateWishlistRequest{
				CountryName: "France",
				Status:      "Planning",
			},
			wantErr:   true,
			errSubstr: "status must be 'Planned' or 'Visited'",
		},
		{
			name: "status with wrong case",
			request: &CreateWishlistRequest{
				CountryName: "Germany",
				Status:      "planned",
			},
			wantErr:   true,
			errSubstr: "status must be 'Planned' or 'Visited'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errSubstr != "" && err.Error() != tt.errSubstr {
				t.Errorf("Validate() error message = %s, want substring %s", err.Error(), tt.errSubstr)
			}
		})
	}
}

func TestUpdateWishlistRequestValidateSuccess(t *testing.T) {
	tests := []struct {
		name    string
		request *UpdateWishlistRequest
		wantErr bool
	}{
		{
			name: "valid request with Planned status",
			request: &UpdateWishlistRequest{
				Note:   "Updated note",
				Status: "Planned",
			},
			wantErr: false,
		},
		{
			name: "valid request with Visited status",
			request: &UpdateWishlistRequest{
				Note:   "Mark as visited",
				Status: "Visited",
			},
			wantErr: false,
		},
		{
			name: "valid request with empty note",
			request: &UpdateWishlistRequest{
				Note:   "",
				Status: "Planned",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateWishlistRequestValidateErrors(t *testing.T) {
	tests := []struct {
		name      string
		request   *UpdateWishlistRequest
		wantErr   bool
		errSubstr string
	}{
		{
			name: "missing status",
			request: &UpdateWishlistRequest{
				Note:   "Some note",
				Status: "",
			},
			wantErr:   true,
			errSubstr: "status is required",
		},
		{
			name: "invalid status",
			request: &UpdateWishlistRequest{
				Note:   "Some note",
				Status: "InProgress",
			},
			wantErr:   true,
			errSubstr: "status must be 'Planned' or 'Visited'",
		},
		{
			name: "status with wrong case",
			request: &UpdateWishlistRequest{
				Note:   "Some note",
				Status: "visited",
			},
			wantErr:   true,
			errSubstr: "status must be 'Planned' or 'Visited'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errSubstr != "" && err.Error() != tt.errSubstr {
				t.Errorf("Validate() error message = %s, want substring %s", err.Error(), tt.errSubstr)
			}
		})
	}
}

func TestWishlistEntryNewInstance(t *testing.T) {
	testDate := time.Date(2023, time.June, 15, 0, 0, 0, 0, time.UTC)
	w := &WishlistEntry{
		ID:          "entry-123",
		CountryName: "France",
		Slug:        "france",
		Note:        "Must visit Paris",
		Status:      StatusPlanned,
		CreatedAt:   testDate,
	}

	if w.ID != "entry-123" {
		t.Errorf("ID = %s, want entry-123", w.ID)
	}
	if w.CountryName != "France" {
		t.Errorf("CountryName = %s, want France", w.CountryName)
	}
	if w.Slug != "france" {
		t.Errorf("Slug = %s, want france", w.Slug)
	}
	if w.Status != StatusPlanned {
		t.Errorf("Status = %s, want %s", w.Status, StatusPlanned)
	}
}

func TestWishlistStatusConstants(t *testing.T) {
	if StatusPlanned != "Planned" {
		t.Errorf("StatusPlanned = %s, want Planned", StatusPlanned)
	}
	if StatusVisited != "Visited" {
		t.Errorf("StatusVisited = %s, want Visited", StatusVisited)
	}
}

func TestValidStatusesMap(t *testing.T) {
	if !ValidStatuses[StatusPlanned] {
		t.Errorf("StatusPlanned should be in ValidStatuses")
	}
	if !ValidStatuses[StatusVisited] {
		t.Errorf("StatusVisited should be in ValidStatuses")
	}
	if ValidStatuses["Invalid"] {
		t.Errorf("Invalid status should not be in ValidStatuses")
	}
}
