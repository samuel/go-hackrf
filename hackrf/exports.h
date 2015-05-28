#ifndef _EXPORTS_H_
#define _EXPORTS_H_ 1

#include <libhackrf/hackrf.h>

extern const hackrf_sample_block_cb_fn *rxCBPtr;
extern const hackrf_sample_block_cb_fn *txCBPtr;

extern int cbGo(hackrf_transfer* transfer, int tx);
int rxCB(hackrf_transfer* transfer);
int txCB(hackrf_transfer* transfer);

#endif
