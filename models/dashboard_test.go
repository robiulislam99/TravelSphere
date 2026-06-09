package models

import (
	"testing"
)

func TestDashboardSummaryNewInstance(t *testing.T) {
	summary := &DashboardSummary{
		Total:   10,
		Planned: 6,
		Visited: 4,
	}

	if summary.Total != 10 {
		t.Errorf("Total = %d, want 10", summary.Total)
	}
	if summary.Planned != 6 {
		t.Errorf("Planned = %d, want 6", summary.Planned)
	}
	if summary.Visited != 4 {
		t.Errorf("Visited = %d, want 4", summary.Visited)
	}
}

func TestDashboardSummaryZeroValues(t *testing.T) {
	summary := &DashboardSummary{}

	if summary.Total != 0 {
		t.Errorf("Total = %d, want 0", summary.Total)
	}
	if summary.Planned != 0 {
		t.Errorf("Planned = %d, want 0", summary.Planned)
	}
	if summary.Visited != 0 {
		t.Errorf("Visited = %d, want 0", summary.Visited)
	}
}

func TestDashboardSummaryAllPlanned(t *testing.T) {
	summary := &DashboardSummary{
		Total:   5,
		Planned: 5,
		Visited: 0,
	}

	if summary.Total != 5 {
		t.Errorf("Total = %d, want 5", summary.Total)
	}
	if summary.Planned != 5 {
		t.Errorf("Planned = %d, want 5", summary.Planned)
	}
	if summary.Visited != 0 {
		t.Errorf("Visited = %d, want 0", summary.Visited)
	}
}

func TestDashboardSummaryAllVisited(t *testing.T) {
	summary := &DashboardSummary{
		Total:   8,
		Planned: 0,
		Visited: 8,
	}

	if summary.Total != 8 {
		t.Errorf("Total = %d, want 8", summary.Total)
	}
	if summary.Planned != 0 {
		t.Errorf("Planned = %d, want 0", summary.Planned)
	}
	if summary.Visited != 8 {
		t.Errorf("Visited = %d, want 8", summary.Visited)
	}
}

func TestDashboardSummaryLargeNumbers(t *testing.T) {
	summary := &DashboardSummary{
		Total:   1000,
		Planned: 600,
		Visited: 400,
	}

	if summary.Total != 1000 {
		t.Errorf("Total = %d, want 1000", summary.Total)
	}
	if summary.Planned != 600 {
		t.Errorf("Planned = %d, want 600", summary.Planned)
	}
	if summary.Visited != 400 {
		t.Errorf("Visited = %d, want 400", summary.Visited)
	}
}
