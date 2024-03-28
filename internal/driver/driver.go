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
const OUT_ENDPOINT = 0x01
const IN_ENDPOINT = 0x82

func StartDriver(configurations *config.Configurations) {

	device, err := util.FindHwMonDevice(&configurations.CPUDevices)
	if err != nil {
		log.Fatalf("Unable to initialize driver: %v", err)
	}

	for {
		err = driverLoop(device, configurations)
		if err != nil {
			log.Printf("Error while communicating with the device: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

}

func driverLoop(device *config.HhMonDeviceRegistered, configurations *config.Configurations) error {
	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(VID, PID)
	if err != nil {
		log.Println("Could not open a device")
		return err

	}
	defer dev.Close()
	dev.SetAutoDetach(true)

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Printf("%s.DefaultInterface()\n", dev)
		return err
	}
	defer done()

	// Open an OUT endpoint.
	ep, err := intf.OutEndpoint(OUT_ENDPOINT)
	if err != nil {
		log.Printf("%s.OutEndpoint(0x01)\n", intf)
		return err
	}

	// Open an IN endpoint
	in, err := intf.InEndpoint(IN_ENDPOINT)
	if err != nil {
		log.Printf("%s.InEndpoint(0x82)\n", intf)
		return err
	}

	for {

		cpuTemp := util.GetCpuTemp(device.GetTemperatureMonitorFile(), configurations.TemperatureScale)
		cpuUsage := util.GetCpuUsage()
		data := []byte(fmt.Sprintf("HLXDATA(%f,%v,0,0,%s)\r\n", cpuUsage, cpuTemp, configurations.TemperatureScale))
		data = append(data, make([]byte, 0x40-len(data))...)

		// Write data to the USB device.
		numBytes, err := ep.Write(data)
		if numBytes != len(data) || err != nil {
			fmt.Printf("%s.Write([5]): only %d bytes written.\n", ep, numBytes)
			return err
		}

		// Read data to the USB device.
		res := make([]byte, 0x40)
		numBytes, err = in.Read(res)

		if numBytes != len(res) || err != nil {
			fmt.Printf("%s.Read([5]): only %d bytes readed\n", ep, numBytes)
			return err
		}

		stringRes := string(res)

		if !strings.HasPrefix(stringRes, "OK(02)") {
			fmt.Printf("Device error message: %s", stringRes)
		}

		time.Sleep(configurations.IntervalTime * time.Second)

	}
}
