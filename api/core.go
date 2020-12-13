package api

import (
	"os"
	"runtime"
	"time"

	"github.com/bodenr/opsyc/util"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type RuntimeEnv struct {
	Timestamp        time.Time
	TimeZone         string
	GoVersion        string
	Hostname         string
	OperatingSystem  string
	Arch             string
	VirtualMemory    *mem.VirtualMemoryStat
	ProcessorPercent []float64
	Interfaces       []net.InterfaceStat
}

func Hostname() string {
	host, err := os.Hostname()
	if err != nil {
		util.Log.Error().Msg("unable to determine hostname: " + err.Error())
		host = "unknown"
	}
	return host
}

func NewRuntimeEnv() *RuntimeEnv {

	virtMem, err := mem.VirtualMemory()
	if err != nil {
		util.Log.Error().Msg("unable to determine virtual memory: " + err.Error())
	}
	cpu, err := cpu.Percent(0, true)
	if err != nil {
		util.Log.Error().Msg("unable to determine CPU: " + err.Error())
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		util.Log.Error().Msg("unable to determine network interfaces: " + err.Error())
	}

	now := time.Now()
	zone, _ := now.Zone()

	return &RuntimeEnv{
		Timestamp:        now,
		TimeZone:         zone,
		Hostname:         Hostname(),
		GoVersion:        runtime.Version(),
		OperatingSystem:  runtime.GOOS,
		Arch:             runtime.GOARCH,
		VirtualMemory:    virtMem,
		ProcessorPercent: cpu,
		Interfaces:       ifaces,
	}
}
