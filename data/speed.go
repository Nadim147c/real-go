package data

import (
	"fmt"
	"time"
)

// Speed represents a quantity of data transfer in bytes per second.
type Speed uint64

// NewSpeed creates a speed from given amount and time.
func NewSpeed(amount Size, dur time.Duration) Speed {
	if dur < 0 {
		panic(fmt.Sprintf("negative time found: %d", dur))
	}

	// Calculate bytes per second
	if dur == 0 {
		return Speed(0)
	}
	return Speed(amount*Size(time.Second)) / Speed(dur)
}

// Size returns the speed as a Size (bytes per second)
func (s Speed) Size() Size {
	return Size(s)
}

// FormatUnitString formats the Speed using the specified unit and precision.
//
// Supported units include all units supported by Size.FormatUnitString,
// with "/s" appended for per-second notation.
func (s Speed) FormatUnitString(unit string, precision ...int) string {
	if s == 0 {
		return "0 " + unit + "/s"
	}

	// Convert to Size and use its FormatUnitString, then append "/s"
	size := s.Size()
	formatted := size.FormatUnitString(unit, precision...)

	// Append "/s" to the formatted string
	// The formatted string is in format "value unit" or "value.fraction unit"
	// We need to insert "/s" before the space or at the end
	return formatted + "/s"
}

// Format implements fmt.Formatter. Supported verbs:
//   - %B for binary byte units per second (KiB/s, MiB/s, ...)
//   - %b for binary bit units per second (Kib/s, Mib/s, ...)
//   - %M for metric byte units per second (kB/s, MB/s, ...)
//   - %m for metric bit units per second (Kb/s, Mb/s, ...)
//   - %d for the raw uint64 value
//   - %s for a string representation similar to %B but ignoring precision
func (s Speed) Format(f fmt.State, verb rune) {
	s.Size().Format(f, verb)
	fmt.Fprint(f, "/s")
}

// String returns the default string representation of the Speed.
//
// It uses binary byte units per second and prints with two decimal places,
// except for raw bytes per second, which are printed as integers.
func (s Speed) String() string {
	return s.Size().String() + "/s"
}

// BytesPerSecond returns the speed in bytes per second as a uint64
func (s Speed) BytesPerSecond() uint64 {
	return uint64(s)
}

// KilobitsPerSecond returns the speed in kilobits per second (metric)
func (s Speed) KilobitsPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(Kb)
}

// MegabitsPerSecond returns the speed in megabits per second (metric)
func (s Speed) MegabitsPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(Mb)
}

// KilobytesPerSecond returns the speed in kilobytes per second (metric)
func (s Speed) KilobytesPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(KB)
}

// MegabytesPerSecond returns the speed in megabytes per second (metric)
func (s Speed) MegabytesPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(MB)
}

// KibibitsPerSecond returns the speed in kibibits per second (binary)
func (s Speed) KibibitsPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(Kib)
}

// MebibitsPerSecond returns the speed in mebibits per second (binary)
func (s Speed) MebibitsPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(Mib)
}

// KibibytesPerSecond returns the speed in kibibytes per second (binary)
func (s Speed) KibibytesPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(KiB)
}

// MebibytesPerSecond returns the speed in mebibytes per second (binary)
func (s Speed) MebibytesPerSecond() float64 {
	if s == 0 {
		return 0
	}
	size := s.Size()
	return float64(size) / float64(MiB)
}
