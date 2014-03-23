package driver 
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func IoInit() bool {
	return int(C.io_init()) != 1
}

func SetBit(channel int) {
	C.io_set_bit(C.int(channel))
}

func ClearBit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func ReadBit(channel int) bool {
	return int(C.io_read_bit(C.int(channel))) != 0
}

func WriteAnalog(channel, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

func ReadAnalog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}

