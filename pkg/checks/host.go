package checks

import (
	"fmt"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/framework"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/hako/durafmt"
	"github.com/shirou/gopsutil/v3/host"
)

var getHostInfo func() (*host.InfoStat, error) = host.Info

// Host collects inspection information on the host machine's metadata, such
// as the operating system
type Host struct{}

// Describe provides a textual description of what this check gathers
func (*Host) Describe() string {
	return "operating system"
}

// Run executes the Host inspection checks
func (host *Host) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	go func() {
		hostInfo, err := getHostInfo()
		if err != nil {
			log.Debug("Unable to inspect host info: %s", err)
			future <- []framework.CheckResult{
				{
					Title:  "Error",
					Status: framework.STATUS_ERROR,
					Value:  fmt.Sprintf("%s", err),
				},
			}

			return
		}

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
