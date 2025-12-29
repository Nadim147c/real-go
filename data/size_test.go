package data

import (
	"fmt"
	"testing"
)

func TestSize_quotient(t *testing.T) {
	tests := []struct {
		name string
		d    Size
		u    Size
		want float64
	}{
		{"exact division", 1024, 1024, 1.0},
		{"fractional division", 1536, 1024, 1.5},
		{"smaller than unit", 512, 1024, 0.5},
		{"zero unit", 1, 0, 0}, // NaN check below
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.quotient(tt.u)
			if tt.u == 0 {
				if got == got {
					t.Fatalf("expected NaN, got %v", got)
				}
				return
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatUnitString(t *testing.T) {
	tests := []struct {
		name      string
		size      Size
		precision int
		unit      string
		want      string
	}{
		{"zero bytes", 0, 0, "B", "0 B"},
		{"raw bytes no precision", 42, 0, "B", "42 B"},
		{"raw bytes with precision", 42, 2, "B", "42.00 B"},
		{"raw bits", 1, 0, "b", "8 b"},
		{"raw bits with precision", 1, 3, "b", "8.000 b"},
		{"metric bytes", 1500, 2, "kB", "1.50 kB"},
		{"binary bytes", 1536, 2, "KiB", "1.50 KiB"},
		{"metric bits", 1000, 2, "kb", "8.00 kb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.size.FormatUnitString(tt.unit, tt.precision)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBestUnit(t *testing.T) {
	tests := []struct {
		name string
		size Size
		unit FormatUnit
		want string
	}{
		{"bytes stay bytes", 512, FormatBinaryByte, "B"},
		{"binary KiB", 1024, FormatBinaryByte, "kiB"},
		{"binary MiB", 5 * MiB, FormatBinaryByte, "MiB"},
		{"metric KB", 1000, FormatMetricByte, "kB"},
		{"metric bits", 1000, FormatMetricBit, "kb"},
		{"binary bits", 1024, FormatBinaryBit, "kib"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.size.bestUnit(tt.unit)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		size Size
		want string
	}{
		{"zero", 0, "0 B"},
		{"bytes", 42, "42 B"},
		{"KiB formatting", 1536, "1.50 kiB"},
		{"MiB formatting", 5 * MiB, "5.00 MiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.size.String()
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFmtFormatting(t *testing.T) {
	tests := []struct {
		name string
		fmt  string
		size Size
		want string
	}{
		{"binary byte verb", "%B", 1536, "1.50 kiB"},
		{"metric byte verb", "%M", 1500, "1.50 kB"},
		{"binary bit verb", "%b", 1024, "8.00 kib"},
		{"metric bit verb", "%m", 1000, "8.00 kb"},
		{"raw int", "%d", 1234, "1234"},
		{"precision override", "%.1B", 1536, "1.5 kiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf(tt.fmt, tt.size)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
