package real

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
		{"celsius", Celsius(0), TempFreezing},
		{"fahrenheit", Fahrenheit(32), TempFreezing},
		{"fahrenheit boiling", Fahrenheit(212), TempBoiling},
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
		unit TempUnit
		want float64
	}{
		{"kelvin to kelvin", TempFreezing, TempUnitKelvin, 273.15},
		{"kelvin to celsius", TempFreezing, TempUnitCelsius, 0},
		{"kelvin to fahrenheit", TempFreezing, TempUnitFahrenheit, 32},
		{"boiling to celsius", TempBoiling, TempUnitCelsius, 100},
		{"boiling to fahrenheit", TempBoiling, TempUnitFahrenheit, 212},
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

func TestTemperature_String(t *testing.T) {
	tests := []struct {
		name string
		t    Temperature
		want string
	}{
		{"zero", 0, "-273.15 °C"},
		{"freezing", TempFreezing, "0.00 °C"},
		{"boiling", TempBoiling, "100.00 °C"},
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
		{"kelvin", "%K", TempFreezing, "273.15 K"},
		{"celsius", "%C", TempFreezing, "0.00 °C"},
		{"fahrenheit", "%F", TempFreezing, "32.00 °F"},
		{"alias f", "%f", TempFreezing, "0.00 °C"},
		{"precision override", "%.1C", TempFreezing, "0.0 °C"},
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

	_ = TempFreezing.In(TempUnit(999))
}
