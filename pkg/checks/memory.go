package checks

import (
	"fmt"

	"github.com/cyberark/conjur-inspect/pkg/framework"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/mem"
)

var getVirtualMemory func() (*mem.VirtualMemoryStat, error) = mem.VirtualMemory

// Memory collects inspection information on the host machine's memory
// availability and usage
type Memory struct{}

// Describe provides a textual description of what this check gathers info on
func (*Memory) Describe() string {
	return "memory"
}

// Run executes the Memory inspection checks
func (memory *Memory) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	go func() {
		v, err := getVirtualMemory()
		if err != nil {
			log.Debug("Unable to inspect memory: %s", err)
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
			{
				Title:  "Memory Total",
				Status: framework.STATUS_INFO,
				Value:  humanize.Bytes(v.Total),
			},
			{
				Title:  "Memory Free",
				Status: framework.STATUS_INFO,
				Value:  humanize.Bytes(v.Free),
			},
			{
				Title:  "Memory Used",
				Status: framework.STATUS_INFO,
				Value: fmt.Sprintf(
					"%s (%.1f %%)",
					humanize.Bytes(v.Used),
					v.UsedPercent,
				),
			},
		}
	}() // async

	return future
}
