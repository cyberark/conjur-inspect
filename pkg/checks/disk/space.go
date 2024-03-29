package disk

import (
	"fmt"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
)

// Aliasing our external dependencies like this allows to swap them out for
// testing
var getPartitions func(all bool) ([]disk.PartitionStat, error) = disk.Partitions
var getUsage func(path string) (*disk.UsageStat, error) = disk.Usage

// SpaceCheck reports on the available partitions and devices on the current
// machine, as well as their available disk space.
type SpaceCheck struct{}

// Describe provides a textual description of what this check gathers info on
func (*SpaceCheck) Describe() string {
	return "disk capacity"
}

// Run executes the disk checks and returns their results
func (*SpaceCheck) Run(context *check.RunContext) <-chan []check.Result {
	future := make(chan []check.Result)

	go func() {
		partitions, err := getPartitions(true)
		// If we can't list the partitions, we exit early with the failure message
		if err != nil {
			log.Debug("Unable to list disk partitions: %s", err)
			future <- []check.Result{
				{
					Title:  "Error",
					Status: check.StatusError,
					Value:  fmt.Sprintf("%s", err),
				},
			}

			return
		}

		results := []check.Result{}

		for _, partition := range partitions {
			usage, err := getUsage(partition.Mountpoint)
			if err != nil {
				log.Debug(
					"Unable to collect disk usage for '%s': %s",
					partition.Mountpoint,
					err,
				)
				continue
			}

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
) check.Result {
	return check.Result{
		Title: fmt.Sprintf(
			"Disk Space (%s, %s)",
			usage.Fstype,
			partition.Mountpoint,
		),
		Status: check.StatusInfo,
		Value: fmt.Sprintf(
			"%s Total, %s Used (%s), %s Free",
			humanize.Bytes(usage.Total),
			humanize.Bytes(usage.Used),
			fmt.Sprintf("%2.f%%", usage.UsedPercent),
			humanize.Bytes(usage.Free),
		),
	}
}
