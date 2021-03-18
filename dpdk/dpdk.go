package dpdk

/*
#cgo LDFLAGS: -lpthread -dl -rt
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <inttypes.h>
#include <sys/types.h>
#include <sys/queue.h>
#include <netinet/in.h>
#include <setjmp.h>
#include <stdarg.h>
#include <ctype.h>
#include <errno.h>
#include <getopt.h>

#include <rte_config.h>
#include <rte_common.h>
#include <rte_log.h>
#include <rte_memory.h>
#include <rte_memcpy.h>
#include <rte_memzone.h>
#include <rte_eal.h>
#include <rte_per_lcore.h>
#include <rte_launch.h>
#include <rte_atomic.h>
#include <rte_cycles.h>
#include <rte_prefetch.h>
#include <rte_lcore.h>
#include <rte_per_lcore.h>
#include <rte_branch_prediction.h>
#include <rte_interrupts.h>
#include <rte_pci.h>
#include <rte_random.h>
#include <rte_debug.h>
#include <rte_ether.h>
#include <rte_ethdev.h>
#include <rte_ring.h>
#include <rte_mempool.h>
#include <rte_mbuf.h>
*/
import "C"
import (
	"fmt"
)

type DPDK struct {
	pktmbuf_pool *C.rte_mempool
}

const NB_MBUF = 8192

func getMbufSize() *C.u_int32_t {
	size := C.sizeof(C.Struct(C.rte_mbuf))
	headroom := C.u_int32_t(C.RTE_PKTMBUF_HEADROOM)
	return C.u_int32_t(2048 + size + headroom)
}

func createMemPool(mpool_name string) (*C.rte_mempool, error) {
	ptr, err := C.rte_mempool_create(
		C.CString(mpool_name),
		C.u_int32_t(NB_MBUF),
		getMbufSize(),
		C.u_int32_t(32),
		C.sizeof(C.Struct(C.rte_pktmbuf_pool_private)),
		C.rte_pktmbuf_pool_init, nil,
		C.rte_pktmbuf_init, nil,
		C.rte_socket_id(), C.u_int32_t(0),
	)

	return ptr, err
}

// InitDpdk initialize intel dpdk & return instance
func InitDpdk(device string) (dpdk *DPDK, _ error) {
	cptr, err := C.rte_eal_init()
	if cptr == nil || err != nil {
		return nil, fmt.Errorf("dpdk NewDataPlane error: %v", err)
	}

	mbufptr, err := createMemPool("mbuf_pool")
	if mbufptr == nil || err != nil {
		return nil, fmt.Errorf("dpdk error calling rte_mempool_create: %v", err)
	}

	dpdk = &DPDK{pktmbuf_pool: mbufptr}

	return
}
