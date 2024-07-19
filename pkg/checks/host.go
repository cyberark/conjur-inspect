package checks

import (
	"fmt"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/check"
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
func (h *Host) Run(context *check.RunContext) []check.Result {
	hostInfo, err := getHostInfo()
	if err != nil {
		return check.ErrorResult(
			h,
			fmt.Errorf("failed to collect host info: %w", err),
		)
	}

	return []check.Result{
		hostnameResult(hostInfo),
		uptimeResult(hostInfo),
		osResult(hostInfo),
		virtualizationResult(hostInfo),
	}
}

func hostnameResult(hostInfo *host.InfoStat) check.Result {
	return check.Result{
		Title:  "Hostname",
		Status: check.StatusInfo,
		Value:  hostInfo.Hostname,
	}
}

func uptimeResult(hostInfo *host.InfoStat) check.Result {
	return check.Result{
		Title:  "Uptime",
		Status: check.StatusInfo,
		Value:  durafmt.Parse(time.Duration(hostInfo.Uptime) * time.Second).String(),
	}
}

func osResult(hostInfo *host.InfoStat) check.Result {
	return check.Result{
		Title:  "OS",
		Status: check.StatusInfo,
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
			Status: check.StatusInfo,
			Value:  "None",
		}
	}

	return check.Result{
		Title:  "Virtualization",
		Status: check.StatusInfo,
		Value: fmt.Sprintf(
			"%s (%s)",
			hostInfo.VirtualizationSystem,
			hostInfo.VirtualizationRole,
		),
	}
}
