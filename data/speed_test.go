package data

import (
	"testing"
	"time"
)

func TestNewSpeed(t *testing.T) {
	tests := []struct {
		name     string
		amount   Size
		dur      time.Duration
		expected Speed
		panic    bool
	}{
		{
			name:     "1 MB in 1 second",
			amount:   MB,
			dur:      time.Second,
			expected: Speed(MB), // 1 MB/s
		},
		{
			name:     "100 MB in 2 seconds",
			amount:   100 * MB,
			dur:      2 * time.Second,
			expected: Speed(50 * MB), // 50 MB/s
		},
		{
			name:     "1 MiB in 1 second",
			amount:   MiB,
			dur:      time.Second,
			expected: Speed(MiB), // 1 MiB/s
		},
		{
			name:     "1 byte in 1 second",
			amount:   Byte,
			dur:      time.Second,
			expected: Speed(Byte), // 1 B/s
		},
		{
			name:     "0 bytes in 1 second",
			amount:   Zero,
			dur:      time.Second,
			expected: 0,
		},
		{
			name:     "100 MB in 0 seconds (division by zero)",
			amount:   100 * MB,
			dur:      0,
			expected: 0,
		},
		{
			name:     "500 MB in 500ms",
			amount:   500 * MB,
			dur:      500 * time.Millisecond,
			expected: Speed(1000 * MB), // 1000 MB/s
		},
		{
			name:   "negative duration - should panic",
			amount: MB,
			dur:    -time.Second,
			panic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.panic {
					t.Errorf("unexpected panic: %v", r)
				} else if r == nil && tt.panic {
					t.Error("expected panic but didn't get one")
				}
			}()

			got := NewSpeed(tt.amount, tt.dur)
			if got != tt.expected {
				t.Errorf("NewSpeed(%v, %v) = %v, want %v",
					tt.amount, tt.dur, got, tt.expected)
			}
		})
	}
}

func TestParseSpeed(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Speed
		wantErr bool
	}{
		{
			name:  "bytes per second",
			input: "100B/s",
			want:  100,
		},
		{
			name:  "kilobytes per second",
			input: "1KB/s",
			want:  Speed(KB),
		},
		{
			name:  "bytes per millisecond",
			input: "1B/ms",
			want:  Speed(1000),
		},
		{
			name:  "megabytes per second",
			input: "2MB/s",
			want:  Speed(2 * MB),
		},
		{
			name:  "whitespace tolerant",
			input: " 1 KB / s ",
			want:  Speed(KB),
		},
		{
			name:    "negative size allowed",
			input:   "-1KB/s",
			want:    Speed(0), // wraps; design-dependent
			wantErr: true,
		},
		{
			name:    "missing separator",
			input:   "1KBs",
			wantErr: true,
		},
		{
			name:    "invalid duration",
			input:   "1KB/day",
			wantErr: true,
		},
		{
			name:    "invalid size",
			input:   "XB/s",
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSpeed(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("ParseSpeed(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSpeed_FormatUnitString(t *testing.T) {
	tests := []struct {
		name     string
		speed    Speed
		unit     string
		prec     []int
		expected string
	}{
		// Byte units
		{
			name:     "0 B/s",
			speed:    0,
			unit:     "B",
			expected: "0 B/s",
		},
		{
			name:     "1024 B/s",
			speed:    Speed(1024),
			unit:     "B",
			expected: "1024 B/s",
		},
		{
			name:     "1024 B/s with 2 decimal",
			speed:    Speed(1024),
			unit:     "B",
			prec:     []int{2},
			expected: "1024.00 B/s",
		},
		{
			name:     "1024 KiB/s",
			speed:    Speed(KiB),
			unit:     "kiB",
			expected: "1 kiB/s",
		},
		{
			name:     "1.5 MiB/s",
			speed:    Speed(MiB) + Speed(512*KiB), // 1.5 MiB/s
			unit:     "MiB",
			prec:     []int{2},
			expected: "1.50 MiB/s",
		},
		{
			name:     "2.5 GB/s",
			speed:    Speed(2*GB) + Speed(500*MB),
			unit:     "GB",
			prec:     []int{1},
			expected: "2.5 GB/s",
		},

		// Bit units
		{
			name:     "8 b/s",
			speed:    Speed(1), // 1 byte/s = 8 bits/s
			unit:     "b",
			expected: "8 b/s",
		},
		{
			name:     "8192 b/s",
			speed:    Speed(1024), // 1024 bytes/s = 8192 bits/s
			unit:     "b",
			prec:     []int{2},
			expected: "8192.00 b/s",
		},
		{
			name:     "1 Mb/s",
			speed:    Speed(MB) / 8, // 1 MB/s = 8 Mb/s
			unit:     "Mb",
			prec:     []int{2},
			expected: "1.00 Mb/s",
		},
		{
			name:     "1 Mib/s",
			speed:    Speed(MiB),
			unit:     "Mib",
			prec:     []int{2},
			expected: "8.00 Mib/s",
		},

		// Metric units
		{
			name:     "1000 kB/s",
			speed:    Speed(MB),
			unit:     "kB",
			expected: "1000 kB/s",
		},
		{
			name:     "1 MB/s",
			speed:    Speed(MB),
			unit:     "MB",
			expected: "1 MB/s",
		},

		// Binary units
		{
			name:     "1024 kiB/s",
			speed:    Speed(MiB),
			unit:     "kiB",
			expected: "1024 kiB/s",
		},
		{
			name:     "1 GiB/s",
			speed:    Speed(GiB),
			unit:     "GiB",
			prec:     []int{3},
			expected: "1.000 GiB/s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.speed.FormatUnitString(tt.unit, tt.prec...)
			if got != tt.expected {
				t.Errorf("FormatUnitString(%q, %v) = %q, want %q",
					tt.unit, tt.prec, got, tt.expected)
			}
		})
	}
}

func TestSpeed_String(t *testing.T) {
	tests := []struct {
		name     string
		speed    Speed
		expected string
	}{
		{
			name:     "0 speed",
			speed:    0,
			expected: "0 B/s",
		},
		{
			name:     "bytes per second",
			speed:    Speed(999),
			expected: "999 B/s",
		},
		{
			name:     "just under 1 KiB/s",
			speed:    Speed(1023),
			expected: "1023 B/s",
		},
		{
			name:     "exactly 1 KiB/s",
			speed:    Speed(KiB),
			expected: "1.00 kiB/s",
		},
		{
			name:     "1.5 MiB/s",
			speed:    Speed(MiB) + Speed(512*KiB),
			expected: "1.50 MiB/s",
		},
		{
			name:     "1024 MiB/s",
			speed:    Speed(GiB),
			expected: "1.00 GiB/s",
		},
		{
			name:     "1 GiB/s",
			speed:    Speed(GiB) + Speed(1),
			expected: "1.00 GiB/s",
		},
		{
			name:     "2.5 GiB/s",
			speed:    Speed(2*GiB) + Speed(512*MiB),
			expected: "2.50 GiB/s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.speed.String()
			if got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}
