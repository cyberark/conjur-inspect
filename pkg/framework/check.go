package framework

// STATUS_INFO means the result is informational only
const STATUS_INFO = "INFO"

// STATUS_PASS means the result falls within the production operational requirements
const STATUS_PASS = "PASS"

// STATUS_WARNS means that the system is at risk for production operation
const STATUS_WARN = "WARN"

// STATUS_FAILS means the system is unacceptable for production operation
const STATUS_FAIL = "FAIL"

// STATUS_ERROR means the result could not be obtained
const STATUS_ERROR = "ERROR"

// Check represent a single operation (API call, external program execution,
// etc.) that returns one or more result.
type Check interface {
	Run() <-chan []CheckResult
}

type CheckResult struct {
	Title   string `json:"title"`
	Value   string `json:"value"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
