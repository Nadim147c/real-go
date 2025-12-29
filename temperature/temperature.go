package temperature

import (
	"fmt"
	"math"
)

// Temperature is a real-life temperature stored in kelvin.
type Temperature float64

// Unit represents a temperature unit.
type Unit int

// UnitFunc converts a float64 in that unit to Temperature.
type UnitFunc func(float64) Temperature

// revive:disable exported
const (
	UnitKelvin Unit = iota
	UnitCelsius
	UnitFahrenheit
)

func Kelvin(t float64) Temperature  { return Temperature(t) }
func Celsius(t float64) Temperature { return Temperature(t) + Freezing }
func Fahrenheit(t float64) Temperature {
	return Temperature((t-32)*5/9) + Freezing
}

// revive:enable exported

// Physical constants.
const (
	// AbsoluteZero is the absolute zero temperature.
	AbsoluteZero Temperature = 0
	// Freezing is the freezing point of water (0°C).
	Freezing Temperature = 273.15
	// Boiling is the boiling point of water (100°C).
	Boiling Temperature = 373.15
)

// In converts temperature to the requested unit.
func (t Temperature) In(u Unit) float64 {
	switch u {
	case UnitKelvin:
		return float64(t)
	case UnitCelsius:
		return float64(t - Freezing)
	case UnitFahrenheit:
		return float64(t-Freezing)*9/5 + 32
	default:
		panic("invalid temperature unit")
	}
}

// String returns a human-friendly representation (°C by default).
func (t Temperature) String() string {
	if math.IsNaN(float64(t)) {
		return "0"
	}
	return fmt.Sprintf("%.2C", t)
}

// Format implements fmt.Formatter.
//
// Supported verbs:
//   - %K — kelvin
//   - %C — celsius
//   - %F — fahrenheit
//   - %f — alias for %C
func (t Temperature) Format(f fmt.State, verb rune) {
	precision, ok := f.Precision()
	if !ok {
		precision = 2
	}

	format := fmt.Sprintf("%%.%df %%s", precision)

	switch verb {
	case 'K':
		fmt.Fprintf(f, format, t.In(UnitKelvin), "K")
	case 'C', 'f':
		fmt.Fprintf(f, format, t.In(UnitCelsius), "°C")
	case 'F':
		fmt.Fprintf(f, format, t.In(UnitFahrenheit), "°F")
	default:
		fmt.Fprint(f, t.String())
	}
}
