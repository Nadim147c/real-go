// Package datasize provides utilities for representing, formatting, and
// printing data sizes in metric and binary units, for both bytes and bits.
//
// It supports automatic unit selection, custom precision, and integration with
// the fmt package via custom formatting verbs.
package data

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unicode"

	islices "github.com/Nadim147c/real-go/internal/slices"
)

// FormatUnit describes a family of units used for formatting data sizes.
type FormatUnit int

const (
	// FormatBinaryByte represents binary byte units (KiB, MiB, GiB, ...).
	FormatBinaryByte FormatUnit = iota
	// FormatMetricByte represents metric byte units (kB, MB, GB, ...).
	FormatMetricByte
	// FormatBinaryBit represents binary bit units (Kib, Mib, Gib, ...).
	FormatBinaryBit
	// FormatMetricBit represents metric bit units (Kb, Mb, Gb, ...).
	FormatMetricBit
)

// Size represents a quantity of data in bytes.
//
// It is defined as an int64 and can represent both byte- and bit-based
// quantities through conversion.
type Size int64

// revive:disable exported

const (
	// Zero represents a data size of zero bytes.
	Zero Size = 0

	// Byte represents a single byte.
	Byte Size = 1

	// Metric byte units.
	KB Size = 1000 * Byte
	MB Size = 1000 * KB
	GB Size = 1000 * MB
	TB Size = 1000 * GB
	PB Size = 1000 * TB
	EB Size = 1000 * PB

	// Binary byte units.
	KiB Size = 1024 * Byte
	MiB Size = 1024 * KiB
	GiB Size = 1024 * MiB
	TiB Size = 1024 * GiB
	PiB Size = 1024 * TiB
	EiB Size = 1024 * PiB

	// Metric bit units.
	Kb Size = KB / 8
	Mb Size = MB / 8
	Gb Size = GB / 8
	Tb Size = TB / 8
	Pb Size = PB / 8
	Eb Size = EB / 8

	// Binary bit units.
	Kib Size = KiB / 8
	Mib Size = MiB / 8
	Gib Size = GiB / 8
	Tib Size = TiB / 8
	Pib Size = PiB / 8
)

// revive:enable exported

// ParseSize parses a datasize to Size
func ParseSize(s string) (Size, error) {
	trimmed := strings.TrimSpace(s)
	numEnd := strings.LastIndexFunc(trimmed, unicode.IsDigit) + 1
	if numEnd <= 0 {
		return 0, fmt.Errorf("invalid size format: %q", s)
	}
	num, inputUnit := trimmed[:numEnd], trimmed[numEnd:]
	size, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0, err
	}

	unit := strings.TrimSpace(inputUnit)
	if unit == "" {
		unit = "B" // default unit is byte
	}

	// we want convert mib or tib but not weird mIb
	if all(unit, unicode.IsLower) {
		unit = strings.Map(func(r rune) rune {
			if r == 'i' {
				return r
			}
			return unicode.ToUpper(r)
		}, unit)
	}

	mul, ok := UnitTable[unit]
	if !ok {
		return 0, fmt.Errorf("invalid input unit: %q", inputUnit)
	}

	if size > 0 && size > math.MaxInt64/int64(mul) {
		return 0, fmt.Errorf("size overflows int64: %q", s)
	}
	if size < 0 && size < math.MinInt64/int64(mul) {
		return 0, fmt.Errorf("size overflows int64: %q", s)
	}

	return Size(size) * mul, nil
}

func all(s string, f func(rune) bool) bool {
	for _, r := range s {
		if !f(r) {
			return false
		}
	}
	return true
}

// quotient returns d divided by u as a floating-point value. If u is zero,
// return NaN.
func (d Size) quotient(u Size) float64 {
	if u == 0 {
		return math.NaN()
	}
	abs := d / u
	mod := d % u
	return float64(abs) + float64(mod)/float64(u)
}

// Value returns the underlying int64 value
func (d Size) Value() int64 {
	return int64(d)
}

// UnitTable maps supported unit strings to their corresponding Size values.
var UnitTable = map[string]Size{
	"B":  Byte,
	"kB": KB, "KB": KB, "MB": MB, "GB": GB, "TB": TB, "PB": PB, "EB": EB,
	"kiB": KiB, "KiB": KiB, "MiB": MiB, "GiB": GiB, "TiB": TiB, "PiB": PiB, "EiB": EiB,
	"kb": Kb, "Kb": Kb, "Mb": Mb, "Gb": Gb, "Tb": Tb, "Pb": Pb, "Eb": Eb,
	"kib": Kib, "Kib": Kib, "Mib": Mib, "Gib": Gib, "Tib": Tib, "Pib": Pib,
}

// FormatUnitString formats the Size using the specified unit and precision.
//
// Supported units include:
//   - b, B
//   - kB, KB, MB, GB, TB, PB, EB
//   - kiB, KiB, MiB, GiB, TiB, PiB, EiB
//   - kb, Kb, Mb, Gb, Tb, Pb, Eb
//   - kib, Kib, Mib, Gib, Tib, Pib, Eib
//
// A precision of zero prints an integer value. For bits and bytes, precision
// greater than zero appends a fractional part of zeros.
func (d Size) FormatUnitString(unit string, precision ...int) string {
	if d == 0 {
		return "0 " + unit
	}

	prec := islices.OptionalValue(0, precision)

	// Handle bytes.
	if unit == "B" {
		if prec == 0 {
			return fmt.Sprintf("%d %s", int64(d), unit)
		}
		return fmt.Sprintf("%d.%0*d %s", int64(d), prec, 0, unit)
	}

	// Handle bits.
	if unit == "b" {
		bits := big.NewInt(int64(d))
		bits.Mul(bits, big.NewInt(8))

		if prec == 0 {
			return fmt.Sprintf("%s %s", bits, unit)
		}
		return fmt.Sprintf("%s.%0*d %s", bits, prec, 0, unit)
	}

	u, ok := UnitTable[unit]
	if !ok {
		panic("illegal diskspace unit")
	}

	format := fmt.Sprintf("%%.%df %s", prec, unit)
	return fmt.Sprintf(format, d.quotient(u))
}

// Format implements fmt.Formatter. Supported verbs:
//   - %B for binary byte units (KiB, MiB, ...)
//   - %b for binary bit units (Kib, Mib, ...)
//   - %M for metric byte units (KB, MB, ...)
//   - %m for metric bit units (Kb, Mb, ...)
//   - %d for the raw int64 value
//   - %s for a string representation similar to %B but ignoring precision
func (d Size) Format(f fmt.State, verb rune) {
	precision, fixed := f.Precision()
	var unit string

	switch verb {
	case 'B':
		unit = d.bestUnit(FormatBinaryByte)
	case 'b':
		unit = d.bestUnit(FormatBinaryBit)
	case 'M':
		unit = d.bestUnit(FormatMetricByte)
	case 'm':
		unit = d.bestUnit(FormatMetricBit)
	case 'd':
		fmt.Fprint(f, int64(d))
		return
	default:
		fmt.Fprint(f, d.String())
		return
	}

	if fixed {
		fmt.Fprint(f, d.FormatUnitString(unit, precision))
		return
	}

	if unit == "B" || unit == "b" {
		fmt.Fprint(f, d.FormatUnitString(unit))
		return
	}

	fmt.Fprint(f, d.FormatUnitString(unit, 2))
}

// String returns the default string representation of the Size.
//
// It uses binary byte units and prints with two decimal places, except for raw
// bytes, which are printed as integers.
func (d Size) String() string {
	unit := d.bestUnit(FormatBinaryByte)
	switch unit {
	case "b", "B":
		return d.FormatUnitString(unit)
	default:
		return d.FormatUnitString(unit, 2)
	}
}

type pair struct {
	name  string
	value Size
}

var (
	metricBytes = []pair{
		{"B", Byte},
		{"kB", KB},
		{"MB", MB},
		{"GB", GB},
		{"TB", TB},
		{"PB", PB},
		{"EB", EB},
	}
	metricBits = []pair{
		{"b", 0},
		{"kb", Kb},
		{"Mb", Mb},
		{"Gb", Gb},
		{"Tb", Tb},
		{"Pb", Pb},
		{"Eb", Eb},
	}
	binaryBytes = []pair{
		{"B", Byte},
		{"kiB", KiB},
		{"MiB", MiB},
		{"GiB", GiB},
		{"TiB", TiB},
		{"PiB", PiB},
		{"EiB", EiB},
	}
	binaryBits = []pair{
		{"b", 0},
		{"kib", Kib},
		{"Mib", Mib},
		{"Gib", Gib},
		{"Tib", Tib},
		{"Pib", Pib},
	}
)

// bestUnit returns the most appropriate unit name for the Size within the given
// unit family.
//
// The returned unit is chosen such that the formatted value is less than the
// next larger unit.
func (d Size) bestUnit(u FormatUnit) string {
	var unitList []pair

	switch u {
	case FormatBinaryByte:
		unitList = binaryBytes
	case FormatMetricByte:
		unitList = metricBytes
	case FormatBinaryBit:
		unitList = binaryBits
	case FormatMetricBit:
		unitList = metricBits
	default:
		panic("invalid unit kind")
	}

	p := islices.LastItemFunc(unitList, func(a pair) bool {
		return a.value <= d
	})

	return p.name
}
