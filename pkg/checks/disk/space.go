package disk

import (
	"fmt"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
)

// SpaceCheck reports on the available partitions and devices on the current
// machine, as well as their available disk space.
type SpaceCheck struct {
}

// Describe provides a textual description of what this check gathers info on
func (*SpaceCheck) Describe() string {
	return "disk capacity"
}

// Run executes the disk checks and returns their results
func (*SpaceCheck) Run() <-chan []framework.CheckResult {
	future := make(chan []framework.CheckResult)

	go func() {
		partitions, _ := disk.Partitions(true)
		results := []framework.CheckResult{}

		for _, partition := range partitions {
			usage, _ := disk.Usage(partition.Mountpoint)
			if usage.Total == 0 {
				continue
			}

			results = append(
				results,
				partitionDiskSpaceResult(partition, usage),
			)
		}

		future <- results
	}() // async

	return future
}

func partitionDiskSpaceResult(
	partition disk.PartitionStat,
	usage *disk.UsageStat,
) framework.CheckResult {
	return framework.CheckResult{
		Title: fmt.Sprintf(
			"Disk Space (%s, %s)",
			usage.Fstype,
			partition.Mountpoint,
		),
		Status: framework.STATUS_INFO,
		Value: fmt.Sprintf(
			"%s Total, %s Used (%s), %s Free",
			humanize.Bytes(usage.Total),
			humanize.Bytes(usage.Used),
			fmt.Sprintf("%2.f%%", usage.UsedPercent),
			humanize.Bytes(usage.Free),
		),
	}
}
