package driver

import (
	"fmt"
	"log"
	"opencool/internal/config"
	"opencool/internal/util"
	"strings"
	"time"

	"github.com/google/gousb"
)

const VID = 0x34d3
const PID = 0x1100

func StartDriver(configurations *config.Configurations) {

	device, err := util.FindHwMonDevice(&configurations.CPUDevices)
	if err != nil {
		log.Fatalf("Unable to initialize driver: %v", err)
	}

	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(VID, PID)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	defer dev.Close()
	dev.SetAutoDetach(true)

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("%s.DefaultInterface(): %v", dev, err)
	}
	defer done()

	// Open an OUT endpoint.
	ep, err := intf.OutEndpoint(0x01)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(0x01): %v", intf, err)
	}

	// Open an IN endpoint
	in, err := intf.InEndpoint(0x82)
	if err != nil {
		log.Fatalf("%s.InEndpoint(0x82): %v", intf, err)
	}

	for {

		cpuTemp := util.GetCpuTemp(device.GetTemperatureMonitorFile())
		cpuUsage := util.GetCpuUsage()
		data := []byte(fmt.Sprintf("HLXDATA(%f,%v,0,0,C)", cpuUsage, cpuTemp))
		data = append(data, 0x0d)
		data = append(data, 0x0a)
		data = append(data, make([]byte, 0x40-len(data))...)

		// Write data to the USB device.
		numBytes, err := ep.Write(data)
		if numBytes != len(data) || err != nil {
			log.Fatalf("%s.Write([5]): only %d bytes written, returned error is %v", ep, numBytes, err)
		}

		// Read data to the USB device.
		res := make([]byte, 0x40)
		numBytes, err = in.Read(res)
		stringRes := string(res)

		if !strings.HasPrefix(stringRes, "OK(02)") {
			fmt.Printf("Device error message: %s", stringRes)
		}

		time.Sleep(configurations.IntervalTime * time.Second)

	}

}
