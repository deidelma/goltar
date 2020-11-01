package dbg

type struct debugStatus {
	active bool
}

var status debug


func init() {
	debug.active := false
}

