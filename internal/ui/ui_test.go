package ui

import "testing"

func TestHumanSize(t *testing.T) {
	tests := map[int64]string{0: "0 B", 1023: "1023 B", 1024: "1.0 KB", 1536: "1.5 KB", 1024 * 1024: "1.0 MB"}
	for size, want := range tests {
		if got := HumanSize(size); got != want {
			t.Errorf("HumanSize(%d) = %q, want %q", size, got, want)
		}
	}
}
