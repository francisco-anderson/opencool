package config

import "time"

type Scale string

const (
	Celsius    Scale = "C"
	Fahrenheit Scale = "F"
)

func (s Scale) IsValid() bool {
	return s == Celsius || s == Fahrenheit
}

type Configurations struct {
	IntervalTime     time.Duration
	TemperatureScale Scale
	CPUDevices       []HwMonDevice
	GPUDevices       []HwMonDevice
}

type HwMonDevice struct {
	Name     string
	TempFile string
}
