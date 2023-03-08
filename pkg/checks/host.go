package checks

import (
	"fmt"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/check"
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
func (*Host) Run() <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {
		hostInfo, err := getHostInfo()
		if err != nil {
			log.Debug("Unable to inspect host info: %s", err)
			future <- []check.Result{
				{
					Title:  "Error",
					Status: check.STATUS_ERROR,
					Value:  fmt.Sprintf("%s", err),
				},
			}

			return
		}

		future <- []check.Result{
			hostnameResult(hostInfo),
			uptimeResult(hostInfo),
			osResult(hostInfo),
			virtualizationResult(hostInfo),
		}
	}() // async

	return future
}

func hostnameResult(hostInfo *host.InfoStat) check.Result {
	return check.Result{
		Title:  "Hostname",
		Status: check.STATUS_INFO,
		Value:  hostInfo.Hostname,
	}
}

func uptimeResult(hostInfo *host.InfoStat) check.Result {
	return check.Result{
		Title:  "Uptime",
		Status: check.STATUS_INFO,
		Value:  durafmt.Parse(time.Duration(hostInfo.Uptime) * time.Second).String(),
	}
}

func osResult(hostInfo *host.InfoStat) check.Result {
	return check.Result{
		Title:  "OS",
		Status: check.STATUS_INFO,
		Value: fmt.Sprintf(
			"%s, %s, %s, %s",
			hostInfo.OS,
			hostInfo.Platform,
			hostInfo.PlatformFamily,
			hostInfo.PlatformVersion,
		),
	}
}

func virtualizationResult(hostInfo *host.InfoStat) check.Result {
	if hostInfo.VirtualizationSystem == "" {
		return check.Result{
			Title:  "Virtualization",
			Status: check.STATUS_INFO,
			Value:  "None",
		}
	}

	return check.Result{
		Title:  "Virtualization",
		Status: check.STATUS_INFO,
		Value: fmt.Sprintf(
			"%s (%s)",
			hostInfo.VirtualizationSystem,
			hostInfo.VirtualizationRole,
		),
	}
}
