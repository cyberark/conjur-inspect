package checks

import (
	"fmt"
	"time"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/hako/durafmt"
	"github.com/shirou/gopsutil/v3/host"
)

type Host struct {
}

// Describe provides a textual description of what this check gathers info on
func (*Host) Describe() string {
	return "operating system"
}

func (*Host) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	go func() {
		hostInfo, _ := host.Info()

		future <- []framework.CheckResult{
			hostnameResult(hostInfo),
			uptimeResult(hostInfo),
			osResult(hostInfo),
			virtualizationResult(hostInfo),
		}
	}() // async

	return future
}

func hostnameResult(hostInfo *host.InfoStat) framework.CheckResult {
	return framework.CheckResult{
		Title:  "Hostname",
		Status: framework.STATUS_INFO,
		Value:  hostInfo.Hostname,
	}
}

func uptimeResult(hostInfo *host.InfoStat) framework.CheckResult {
	return framework.CheckResult{
		Title:  "Uptime",
		Status: framework.STATUS_INFO,
		Value:  durafmt.Parse(time.Duration(hostInfo.Uptime) * time.Second).String(),
	}
}

func osResult(hostInfo *host.InfoStat) framework.CheckResult {
	return framework.CheckResult{
		Title:  "OS",
		Status: framework.STATUS_INFO,
		Value: fmt.Sprintf(
			"%s, %s, %s, %s",
			hostInfo.OS,
			hostInfo.Platform,
			hostInfo.PlatformFamily,
			hostInfo.PlatformVersion,
		),
	}
}

func virtualizationResult(hostInfo *host.InfoStat) framework.CheckResult {
	if hostInfo.VirtualizationSystem == "" {
		return framework.CheckResult{
			Title:  "Virtualization",
			Status: framework.STATUS_INFO,
			Value:  "None",
		}
	}

	return framework.CheckResult{
		Title:  "Virtualization",
		Status: framework.STATUS_INFO,
		Value: fmt.Sprintf(
			"%s (%s)",
			hostInfo.VirtualizationSystem,
			hostInfo.VirtualizationRole,
		),
	}
}
