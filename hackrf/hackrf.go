/*
Package hackrf provides an interface to the HackRF SDR hardware.

This package wraps libhackrf using cgo.
*/
package hackrf

// #cgo darwin CFLAGS: -I/usr/local/include
// #cgo darwin LDFLAGS: -L/usr/local/lib
// #cgo LDFLAGS: -lhackrf
// #include <libhackrf/hackrf.h>
import "C"
import (
	"errors"
	"fmt"
)

var (
	ErrInvalidParam        = errors.New("hackrf: invalid param")
	ErrNotFound            = errors.New("hackrf: not found")
	ErrBusy                = errors.New("hackrf: busy")
	ErrNoMem               = errors.New("hackrf: no mem")
	ErrLibUSB              = errors.New("hackrf: libusb error")
	ErrThread              = errors.New("hackrf: thread error")
	ErrStreamingThreadErr  = errors.New("hackrf: streaming thread error")
	ErrStreamingStopped    = errors.New("hackrf: streaming stopped")
	ErrStreamingExitCalled = errors.New("hackrf: streaming exit called")
	ErrOther               = errors.New("hackrf: other error")
)

type ErrUnknown int

func (e ErrUnknown) Error() string {
	return fmt.Sprintf("hackrf: unknown error %d", int(e))
}

// Init must be called once at the start of the program.
func Init() error {
	return toError(C.hackrf_init())
}

// Exit should be called once at the end of the program.
func Exit() error {
	return toError(C.hackrf_exit())
}

func toError(r C.int) error {
	if r == C.HACKRF_SUCCESS {
		return nil
	}
	switch r {
	case C.HACKRF_ERROR_INVALID_PARAM:
		return ErrInvalidParam
	case C.HACKRF_ERROR_NOT_FOUND:
		return ErrNotFound
	case C.HACKRF_ERROR_BUSY:
		return ErrBusy
	case C.HACKRF_ERROR_NO_MEM:
		return ErrNoMem
	case C.HACKRF_ERROR_LIBUSB:
		return ErrLibUSB
	case C.HACKRF_ERROR_THREAD:
		return ErrThread
	case C.HACKRF_ERROR_STREAMING_THREAD_ERR:
		return ErrStreamingThreadErr
	case C.HACKRF_ERROR_STREAMING_STOPPED:
		return ErrStreamingStopped
	case C.HACKRF_ERROR_STREAMING_EXIT_CALLED:
		return ErrStreamingExitCalled
	case C.HACKRF_ERROR_OTHER:
		return ErrOther
	}
	return ErrUnknown(int(r))
}
