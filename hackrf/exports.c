#include "exports.h"

int rxCB(hackrf_transfer* transfer) {
	return cbGo(transfer, 0);
}

int txCB(hackrf_transfer* transfer) {
	return cbGo(transfer, 1);
}

const hackrf_sample_block_cb_fn *rxCBPtr = (hackrf_sample_block_cb_fn*)&rxCB;
const hackrf_sample_block_cb_fn *txCBPtr = (hackrf_sample_block_cb_fn*)&txCB;
