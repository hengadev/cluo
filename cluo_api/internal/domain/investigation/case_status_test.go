package investigation

import "testing"

func TestStatusIsValid(t *testing.T) {
	tests := []struct {
		name  string
		input Status
		want  bool
	}{
		{"in_progress is valid", StatusInProgress, true},
		{"ready is valid", StatusReady, true},
		{"released is valid", StatusReleased, true},
		{"draft is rejected", Status("draft"), false},
		{"empty string is rejected", Status(""), false},
		{"arbitrary string is rejected", Status("pending"), false},
		{"mixed case is rejected", Status("In_Progress"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.IsValid()
			if got != tt.want {
				t.Errorf("Status(%q).IsValid() = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
