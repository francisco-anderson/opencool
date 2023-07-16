package util

import (
	"errors"
	"log"
	"opencool/internal/config"
	"os"
	"strings"
)

const hwmondir = "/sys/class/hwmon"

func FindHwMonDevice(devices *[]config.HwMonDevice) (*config.HhMonDeviceRegistered, error) {

	entries, err := os.ReadDir(hwmondir)
	if err != nil {
		log.Fatal(err)
	}

	devicesRegistered := make([]*config.HhMonDeviceRegistered, 0)

	for _, e := range entries {
		hwMonDeviceDir := hwmondir + "/" + e.Name()
		dat, err := os.ReadFile(hwMonDeviceDir + "/name")
		if err != nil {
			log.Fatal(err)
		}
		r := &config.HhMonDeviceRegistered{
			Name:           strings.ReplaceAll(string(dat), "\n", ""),
			HwMonDeviceDir: hwMonDeviceDir,
		}
		devicesRegistered = append(devicesRegistered, r)
	}

	for _, device := range *devices {
		deviceRegistered := getDeviceRegistered(devicesRegistered, device.Name)
		if deviceRegistered != nil {
			deviceRegistered.Device = &device

			return deviceRegistered, nil
		}
	}

	return nil, errors.New("Unable to find a device supported by the current configuration. Please check the configuration file for your device")

}

func getDeviceRegistered(s []*config.HhMonDeviceRegistered, str string) *config.HhMonDeviceRegistered {
	for _, v := range s {
		if v.Name == str {
			return v
		}
	}
	return nil
}
