package config

type HhMonDeviceRegistered struct {
	Name           string
	HwMonDeviceDir string
	Device         *HwMonDevice
}

func (hwMonDeviceRegistered *HhMonDeviceRegistered) GetTemperatureMonitorFile() *string {
	absolutePath := hwMonDeviceRegistered.HwMonDeviceDir + "/" + hwMonDeviceRegistered.Device.TempFile
	return &absolutePath
}
