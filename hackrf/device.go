package hackrf

// #include <stdlib.h>
// #include <libhackrf/hackrf.h>
// #include "exports.h"
import "C"

import (
	"sync"
	"unsafe"
)

var (
	cbMu      sync.RWMutex
	callbacks []*callbackContext
)

type Device struct {
	cdev      *C.hackrf_device
	callbacks []int
}

type Callback func(buf []byte) error

type callbackContext struct {
	cb  Callback
	tx  bool
	dev *Device
}

//export cbGo
func cbGo(transfer *C.hackrf_transfer, tx C.int) C.int {
	var cbIdx int
	if tx == 0 {
		cbIdx = int(uintptr(transfer.rx_ctx))
	} else {
		cbIdx = int(uintptr(transfer.tx_ctx))
	}
	cbMu.RLock()
	if cbIdx >= len(callbacks) {
		cbMu.RUnlock()
		return -1
	}
	ctx := callbacks[cbIdx]
	cbMu.RUnlock()
	if ctx == nil {
		return -1
	}
	n := int(transfer.valid_length)
	goBuf := (*[1 << 30]byte)(unsafe.Pointer(transfer.buffer))[:n:n]
	if err := ctx.cb(goBuf); err != nil {
		return -1
	}
	return 0
}

// Open returns the HackRF device which provides access to the hardware.
func Open() (*Device, error) {
	var d Device
	if r := C.hackrf_open(&d.cdev); r != C.HACKRF_SUCCESS {
		return nil, toError(r)
	}
	return &d, nil
}

func (d *Device) Close() error {
	e := toError(C.hackrf_close(d.cdev))
	if e == nil {
		d.cdev = nil
	}
	return e
}

func (d *Device) Version() (string, error) {
	ver := (*C.char)(C.malloc(128))
	defer C.free(unsafe.Pointer(ver))
	if r := C.hackrf_version_string_read(d.cdev, ver, 128); r != C.HACKRF_SUCCESS {
		return "", toError(r)
	}
	return C.GoString(ver), nil
}

// StartRX starts sampling sending the IQ samples to the callback.
func (d *Device) StartRX(cb Callback) error {
	cbIx := d.registerCallback(&callbackContext{
		cb:  cb,
		tx:  false,
		dev: d,
	})
	return toError(C.hackrf_start_rx(d.cdev, (*[0]byte)(unsafe.Pointer(C.rxCBPtr)), unsafe.Pointer(uintptr(cbIx))))
}

func (d *Device) StopRX() error {
	return toError(C.hackrf_stop_rx(d.cdev))
}

func (d *Device) StartTX(cb Callback) error {
	cbIx := d.registerCallback(&callbackContext{
		cb:  cb,
		tx:  true,
		dev: d,
	})
	return toError(C.hackrf_start_tx(d.cdev, (*[0]byte)(unsafe.Pointer(C.rxCBPtr)), unsafe.Pointer(uintptr(cbIx))))
}

func (d *Device) StopTX() error {
	return toError(C.hackrf_stop_tx(d.cdev))
}

func (d *Device) registerCallback(ctx *callbackContext) int {
	cbMu.Lock()
	cbIx := -1
	for i, cb := range callbacks {
		if cb == nil {
			cbIx = i
			callbacks[i] = ctx
			break
		}
	}
	if cbIx < 0 {
		cbIx = len(callbacks)
		callbacks = append(callbacks, ctx)
	}
	cbMu.Unlock()
	d.callbacks = append(d.callbacks, cbIx)
	return cbIx
}

func (d *Device) SetFreq(freqHz uint64) error {
	return toError(C.hackrf_set_freq(d.cdev, C.uint64_t(freqHz)))
}

// extern ADDAPI int ADDCALL hackrf_set_freq_explicit(hackrf_device* device,
// 		const uint64_t if_freq_hz, const uint64_t lo_freq_hz,
// 		const enum rf_path_filter path);

// SetSampleRateManual sets the sample rate in hz.
// Preferred rates are 8, 10, 12.5, 16, 20MHz due to less jitter.
func (d *Device) SetSampleRateManual(freqHz, divider int) error {
	return toError(C.hackrf_set_sample_rate_manual(d.cdev, C.uint32_t(freqHz), C.uint32_t(divider)))
}

// SetSampleRate sets the sample rate in hz.
// Preferred rates are 8, 10, 12.5, 16, 20MHz due to less jitter.
func (d *Device) SetSampleRate(freqHz float64) error {
	return toError(C.hackrf_set_sample_rate(d.cdev, C.double(freqHz)))
}

// SetBasebandFilterBandwidth sets the baseband bandwidth.
// Possible values are 1.75/2.5/3.5/5/5.5/6/7/8/9/10/12/14/15/20/24/28MHz.
func (d *Device) SetBasebandFilterBandwidth(hz int) error {
	return toError(C.hackrf_set_baseband_filter_bandwidth(d.cdev, C.uint32_t(hz)))
}

// SetAmpEnable enables or disables the external RX/TX RF amplifier.
func (d *Device) SetAmpEnable(value bool) error {
	var v C.uint8_t
	if value {
		v = 1
	}
	return toError(C.hackrf_set_amp_enable(d.cdev, v))
}

// SetLNAGain sets the gain for the RX low-noise amplifier (IF).
// Range 0-40 step 8db
func (d *Device) SetLNAGain(value int) error {
	return toError(C.hackrf_set_lna_gain(d.cdev, C.uint32_t(value)))
}

// SetVGAGain sets the gain for the RX variable gain amplifier (baseband).
// Range 0-62 step 2db
func (d *Device) SetVGAGain(value int) error {
	return toError(C.hackrf_set_vga_gain(d.cdev, C.uint32_t(value)))
}

// SetTXVGAGain sets the gain for the TX variable gain amplifier (IF).
// Range 0-47 step 1db
func (d *Device) SetTXVGAGain(value int) error {
	return toError(C.hackrf_set_txvga_gain(d.cdev, C.uint32_t(value)))
}

func (d *Device) SetAntennaEnable(enabled bool) error {
	var value C.uint8_t
	if enabled {
		value = 1
	}
	return toError(C.hackrf_set_antenna_enable(d.cdev, C.uint8_t(value)))
}

// ComputeBasebandFilterBWRoundDownLT computse the nearest freq for bw filter
// (manual filter)
func ComputeBasebandFilterBWRoundDownLT(bandwidthHz int) int {
	return int(C.hackrf_compute_baseband_filter_bw_round_down_lt(C.uint32_t(bandwidthHz)))
}

// ComputeBasebandFilterBWRoundDownLT computes the best default value
// depending on sample rate (auto filter)
func ComputeBasebandFilterBW(bandwidthHz int) int {
	return int(C.hackrf_compute_baseband_filter_bw(C.uint32_t(bandwidthHz)))
}
