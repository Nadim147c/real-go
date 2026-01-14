# Real

A collection of real-world values as Go types. For when you need actual units, not
just raw numbers.

## Why?

Because it was fun to write...! Do I add more?

## When to use?

DO NOT USE IT!

## Quick Examples

```go
import (
    "time"
    "github.com/Nadim147c/real-go/data"
    "github.com/Nadim147c/real-go/temperature"
)

// Data with proper units
file := 2*data.GB + 500*data.MB
fmt.Printf("File: %s\n", file)  // "2.50 GB"

// Download speeds that make sense
downloaded := data.NewSpeed(100*data.MB, 2*time.Second)
fmt.Printf("Speed: %s\n", downloaded)  // "50.00 MiB/s"

// Temperatures that don't confuse
room := temperature.Celsius(20)
fmt.Printf("Room: %s\n", room)  // "20.00°C"
```

## What's Inside

### 📦 Data Sizes (`data.Size`)

Bytes, bits, and everything in between. Knows the difference between MB (1000) and MiB (1024).

```go
// Create
file := 2*data.GB + 500*data.MB  // 2.5 GB file
ram := 16 * data.GiB             // 16 GiB RAM

// Print
fmt.Println(file.String())  // "2.50 GB" (auto-chooses unit)
fmt.Printf("%B\n", ram)     // "16.00 GiB" (binary units)
```

### ⚡ Transfer Speeds (`data.Speed`)

Internet speeds, file transfers, downloads. All with proper "/s" units.

```go
// From downloaded amount and time
speed := data.NewSpeed(100*data.MB, 2*time.Second)
fmt.Printf("%s\n", speed)  // "50.00 MiB/s"

// Fancy formatting
fmt.Printf("%.1M\n", speed)  // "50.0 MB/s" (metric, 1 decimal)
fmt.Printf("%m\n", speed)    // "400.00 Mb/s" (bits!)
```

### 🌡️ Temperatures (`temperature`)

Celsius, Fahrenheit, Kelvin. No more guessing which unit you're in.

```go
// Create
room := temperature.Celsius(20)      // 20°C
body := temperature.Fahrenheit(98.6) // 98.6°F
cold := temperature.Kelvin(0)        // 0K (brrr)

// Convert
fmt.Printf("%s = %.1f°F\n", room,
    room.In(temperature.UnitFahrenheit))  // "20.00°C = 68.0°F"
```

## Install

```bash
go get github.com/Nadim147c/real-go
```

# LICENSE

This repository is licensed under [LGPL-3.0](./LICENSE.md).
