package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"opencool/internal/monitor"

	"github.com/alecthomas/kong"
)

var (
	buildTime string
	version   string
)

type versionFlag bool

var CLI struct {
	TempFile     string      `help:"CPU Temperature file" name:"temperature-file" type:"string" required:""`
	IntervalTime int         `help:"Time in seconds to update" name:"monitor-interval" type:"int" required:""`
	Version      versionFlag `help:"Show build version" name:"version" type:"bool"`
}

func main() {
	ctx := kong.Parse(&CLI)
	fmt.Println(ctx.Command())
	monitor.StartMonitoring(&monitor.MonitorParams{TempFile: CLI.TempFile, IntervalTime: time.Second * time.Duration(CLI.IntervalTime)})
}

func (v versionFlag) BeforeApply() error {
	fmt.Printf("Version:\t%s\n", version)
	fmt.Printf("Build time:\t%s\n", buildTime)
	fmt.Printf("OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
	return nil
}
