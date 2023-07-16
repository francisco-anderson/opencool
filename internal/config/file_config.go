package config

import "time"

type Configurations struct {
	IntervalTime time.Duration
	CPUDevices   []HwMonDevice
	GPUDevices   []HwMonDevice
}

type HwMonDevice struct {
	Name     string
	TempFile string
}
