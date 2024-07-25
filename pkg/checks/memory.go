package checks

import (
	"fmt"

	"github.com/cyberark/conjur-inspect/pkg/check"
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
func (memory *Memory) Run(*check.RunContext) []check.Result {
	v, err := getVirtualMemory()
	if err != nil {
		return check.ErrorResult(
			memory,
			fmt.Errorf("failed to inspect memory: %w", err),
		)
	}

	return []check.Result{
		{
			Title:  "Memory Total",
			Status: check.StatusInfo,
			Value:  humanize.Bytes(v.Total),
		},
		{
			Title:  "Memory Free",
			Status: check.StatusInfo,
			Value:  humanize.Bytes(v.Free),
		},
		{
			Title:  "Memory Used",
			Status: check.StatusInfo,
			Value: fmt.Sprintf(
				"%s (%.1f %%)",
				humanize.Bytes(v.Used),
				v.UsedPercent,
			),
		},
	}
}
