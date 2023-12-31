package main

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime"
	"strings"

	"opencool/internal/config"
	"opencool/internal/driver"

	"github.com/alecthomas/kong"
	"github.com/spf13/viper"
)

var (
	buildTime string
	version   string
)

type versionFlag bool

var CLI struct {
	ConfigFile string      `help:"Opencool configuration file" type:"string" required:""`
	Version    versionFlag `help:"Show build version" name:"version" type:"bool"`
}

func main() {
	stdlog.SetOutput(new(LogWriter))

	ctx := kong.Parse(&CLI)
	fmt.Println(ctx.Command())

	viper.SetConfigFile(CLI.ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	var configuration config.Configurations
	err := viper.Unmarshal((&configuration))
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	driver.StartDriver(&configuration)
}

type LogWriter int

func (LogWriter) Write(data []byte) (int, error) {
	logmessage := string(data)
	//fix - disable warnings libusb: interrupted [code -10]
	if strings.Contains(logmessage, "interrupted [code -10]") {
		return len(data), nil
	}
	fmt.Println(logmessage)
	return len(data), nil
}

func (v versionFlag) BeforeApply() error {
	fmt.Printf("Version:\t%s\n", version)
	fmt.Printf("Build time:\t%s\n", buildTime)
	fmt.Printf("OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
	return nil
}
