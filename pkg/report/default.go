package report

import (
	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/checks"
	"github.com/cyberark/conjur-inspect/pkg/checks/disk"
	"github.com/cyberark/conjur-inspect/pkg/container"
)

// NewDefaultReport returns a report containing the standard inspection checks
func NewDefaultReport(
	id string,
	rawDataDir string,
) (Report, error) {

	return NewReport(
		id,
		rawDataDir,
		defaultReportSections(),
	)
}

func defaultReportSections() []Section {
	return []Section{
		// TODO:
		// - Recent load
		{
			Title: "CPU",
			Checks: []check.Check{
				&checks.Cpu{},
			},
		},
		{
			Title: "Disk",
			Checks: []check.Check{
				&disk.SpaceCheck{},
				disk.NewIopsCheck(),
				disk.NewLatencyCheck(),
			},
		},
		{
			Title: "Memory",
			Checks: []check.Check{
				&checks.Memory{},
			},
		},
		// TODO:
		// - ipv6 status
		// {
		// 	Title:  "Network",
		// 	Checks: []framework.Check{},
		// },
		{
			Title: "Host",
			Checks: []check.Check{
				&checks.Host{},
			},
		},
		{
			Title: "Follower",
			Checks: []check.Check{
				&checks.Follower{},
			},
		},
		{
			Title: "Container Runtime",
			Checks: []check.Check{
				&checks.ContainerRuntime{
					Provider: &container.DockerProvider{},
				},
				&checks.ContainerRuntime{
					Provider: &container.PodmanProvider{},
				},
			},
		},
		{
			Title: "Ulimits",
			Checks: []check.Check{
				&checks.Ulimit{},
			},
		},
	}
}
