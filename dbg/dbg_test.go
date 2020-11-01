package dbg

import "testing"

func TestStartStop(t *testing.T) {
	Start()
	Printf("Hello, I'm visible")
	Stop()
	Printf("Hello, I'm not visible")
	Start()
	Printf("Hello, I'm visible")
}
