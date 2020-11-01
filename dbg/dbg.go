package dbg

import (
	"log"
)

type debugStatus struct {
	active bool
}

var debug debugStatus

// Start allows debug statements to printed to stdout
func Start() {
	debug.active = true
}

// Stop prevents further printing of debug statements
func Stop() {
	debug.active = false
}

// Printf shadows log.Printf provided a switched log statement
func Printf(format string, items ...interface{}) {
	if debug.active {
		log.Printf(format, items...)
	}
}

func init() {
	debug.active = false
}
