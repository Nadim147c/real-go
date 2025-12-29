package temperature

import (
	"fmt"
	"math"
	"testing"
)

func TestConstructors(t *testing.T) {
	tests := []struct {
		name string
		got  Temperature
		want Temperature
	}{
		{"kelvin", Kelvin(300), 300},
		{"celsius", Celsius(0), Freezing},
		{"fahrenheit", Fahrenheit(32), Freezing},
		{"fahrenheit boiling", Fahrenheit(212), Boiling},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestIn(t *testing.T) {
	tests := []struct {
		name string
		t    Temperature
		unit Unit
		want float64
	}{
		{"kelvin to kelvin", Freezing, UnitKelvin, 273.15},
		{"kelvin to celsius", Freezing, UnitCelsius, 0},
		{"kelvin to fahrenheit", Freezing, UnitFahrenheit, 32},
		{"boiling to celsius", Boiling, UnitCelsius, 100},
		{"boiling to fahrenheit", Boiling, UnitFahrenheit, 212},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.In(tt.unit)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		t    Temperature
		want string
	}{
		{"zero", 0, "-273.15 °C"},
		{"freezing", Freezing, "0.00 °C"},
		{"boiling", Boiling, "100.00 °C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.String()
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringNaN(t *testing.T) {
	tp := Temperature(math.NaN())
	got := tp.String()
	if got != "0" {
		t.Fatalf("got %q, want %q", got, "0")
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		fmt  string
		t    Temperature
		want string
	}{
		{"kelvin", "%K", Freezing, "273.15 K"},
		{"celsius", "%C", Freezing, "0.00 °C"},
		{"fahrenheit", "%F", Freezing, "32.00 °F"},
		{"alias f", "%f", Freezing, "0.00 °C"},
		{"precision override", "%.1C", Freezing, "0.0 °C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf(tt.fmt, tt.t)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInInvalidUnitPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid unit")
		}
	}()

	_ = Freezing.In(Unit(999))
}
