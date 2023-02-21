package fio

// Result is a data structure that represents the JSON
// output returned from running `fio`.
type Result struct {
	Version string      `json:"fio version"`
	Jobs    []JobResult `json:"jobs"`
}

// JobResult represents the results from an individual fio job. A FioResult
// may include multiple job results.
type JobResult struct {
	Sync  JobModeResult `json:"sync"`
	Read  JobModeResult `json:"read"`
	Write JobModeResult `json:"write"`
}

// JobModeResult represents the measurements for a given test mode
// (e.g. read, write). Not all modes provide all values. The populated
// values depend on the fio job parameters.
type JobModeResult struct {
	Iops       float64 `json:"iops"`
	IopsMin    int64   `json:"iops_min"`
	IopsMax    int64   `json:"iops_max"`
	IopsMean   float64 `json:"iops_mean"`
	IopsStddev float64 `json:"iops_stddev"`

	LatNs ResultStats `json:"lat_ns"`
}

// ResultStats represents the statistical measurements provided by fio.
type ResultStats struct {
	Min        int64      `json:"min"`
	Max        int64      `json:"max"`
	Mean       float64    `json:"mean"`
	StdDev     float64    `json:"stddev"`
	N          int64      `json:"N"`
	Percentile Percentile `json:"percentile"`
}

// Percentile provides a simple interface to return particular statistical
// percentiles from the fio stats results.
type Percentile struct {
	NinetyNinth int64 `json:"99.000000"`
}
