package log

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// We want all logging to go to the standard error stream, since we output data
// on the standard output stream.
var infoLogger = log.New(os.Stderr, "INFO:  ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
var errorLogger = log.New(os.Stderr, "ERROR: ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
var isDebug = false

/*
RecordedError prints an error message to the error log and returns a new error with the given message.
This method can receive also more arguments (e.g an external error) and they will be appended to the given error message.

For example, we have a local method `someMethod()`. This method handles its own error printing and thus we can consume
the error and not append it to the new error message, as follows:

	returnVal, err := someMethod()
	if err != nil {
		return nil, log.RecordedError("failed to run someMethod")
	}

On the other hand, if `someMethod()` is a 3rd party method we want to print also the returned error as it wasn't printed
to the error log. So we'll have the following code:

	returnVal, err := 3rdParty.someMethod()
	if err != nil {
		return nil, log.RecordedError("failed to run someMethod. Reason: %s", err)
	}
*/
func RecordedError(errorMessage string, args ...interface{}) error {
	message := fmt.Sprintf(errorMessage, args...)
	writeLog(errorLogger, "ERROR", message)
	return errors.New(message)
}

// Error prints a message to the error log with the "ERROR" label.
func Error(message string, args ...interface{}) {
	writeLog(errorLogger, "ERROR", message, args...)
}

// Warn prints a message to the info log with the "WARN" label.
func Warn(message string, args ...interface{}) {
	writeLog(infoLogger, "WARN", message, args...)
}

// Info prints a message to the info log with the "INFO" label.
func Info(message string, args ...interface{}) {
	writeLog(infoLogger, "INFO", message, args...)
}

// Debug prints a message to the info log with the "DEBUG" label.
func Debug(infoMessage string, args ...interface{}) {
	if isDebug {
		writeLog(infoLogger, "DEBUG", infoMessage, args...)
	}
}

// EnableDebugMode enables writing DEBUG level messages to the log output.
func EnableDebugMode() {
	isDebug = true
	Debug("Debug mode is enabled")
}

func writeLog(logger *log.Logger, logLevel string, message string, args ...interface{}) {
	// -7 format ensures logs alignment, by padding spaces to log level to ensure 7 characters length.
	// 5 for longest log level, 1 for ':', and a space separator.
	logger.SetPrefix(fmt.Sprintf("%-7s", logLevel+":"))
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	logger.Output(3, message)
}
