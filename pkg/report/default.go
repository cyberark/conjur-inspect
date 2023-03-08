package report

import (
	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/checks"
	"github.com/cyberark/conjur-inspect/pkg/checks/disk"
)

// NewDefaultReport returns a report containing the standard inspection checks
func NewDefaultReport(debug bool) *Report {
	return &Report{
		Sections: []Section{
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
					disk.NewIopsCheck(debug),
					disk.NewLatencyCheck(debug),
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
			// TODO:
			// - Podman version
			// - Docker version
			// {
			// 	Title:  "Container Runtime",
			// 	Checks: []framework.Check{},
			// },
		},
	}
}
