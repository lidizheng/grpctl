package store

import "log"

var verbose bool = false

// SetVerbose sets verbose to true
func SetVerbose() {
	verbose = true
}

// IsVerbose returns if verbose is set to true
func IsVerbose() bool {
	return verbose
}

// Debugf prints a log if verbose is set to true
func Debugf(format string, v ...interface{}) {
	if !IsVerbose() {
		return
	}
	log.Printf(format, v...)
}
