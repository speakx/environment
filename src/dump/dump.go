package dump

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	dumpEventUnknow = iota
	dumpEventTimetick
	dumpEventRequest
	dumpEventResponse
)

var onceRT sync.Once
var dumpCtx *dumpContext

type dumpContext struct {
	dumpFlg           bool
	packetRequest     int64
	packetResponse    int64
	packetRequestAvg  int64
	packetResponseAvg int64
	dumpEventChan     chan int64
	dumpFunc          func(string)
}

// InitDump init
func InitDump(dump bool, dumpFunc func(string)) {
	onceRT.Do(func() {
		dumpCtx = &dumpContext{
			dumpFlg:       dump,
			dumpEventChan: make(chan int64, 100),
			dumpFunc:      dumpFunc,
		}
		dumpCtx.loopRuntimeInfo()
	})
}

// PacketRequestCounter DUMP_EVENT_PacketRequest
func PacketRequestCounter() {
	dumpCtx.dumpEventChan <- dumpEventRequest
}

// PacketResponseCounter DUMP_EVENT_PacketResponse
func PacketResponseCounter() {
	dumpCtx.dumpEventChan <- dumpEventResponse
}

func (d *dumpContext) loopRuntimeInfo() {
	go func() {
		dumpCounter := int64(0)
		dumpPrint := int64(5)
		for {
			e, _ := <-d.dumpEventChan
			switch e {
			case dumpEventTimetick:
				if d.dumpFlg && dumpCounter > 0 && dumpCounter%dumpPrint == 0 {
					if 0 == d.packetRequestAvg {
						d.packetRequestAvg = d.packetRequest / dumpCounter
						d.packetResponseAvg = d.packetResponse / dumpCounter
					} else {
						d.packetRequestAvg += (d.packetRequest / dumpCounter)
						d.packetResponseAvg += (d.packetResponse / dumpCounter)
						d.packetRequestAvg /= 2
						d.packetResponseAvg /= 2
					}
					dumpInfo := fmt.Sprintf("dump rate(cur/avg) Request:%-7v/%-7v Response:%-7v/%-7v goroutine:%v",
						d.packetRequest/dumpCounter, d.packetRequestAvg,
						d.packetResponse/dumpCounter, d.packetResponseAvg,
						runtime.NumGoroutine())
					d.dumpFunc(dumpInfo)
					// cleanup rate
					dumpCounter = 0
					d.packetRequest = 0
					d.packetResponse = 0
				} else {
					dumpCounter++
				}
			case dumpEventRequest:
				d.packetRequest++
			case dumpEventResponse:
				d.packetResponse++
			}
		}
	}()

	go func() {
		for {
			<-time.After(time.Second * 1)
			d.dumpEventChan <- dumpEventTimetick
		}
	}()
}
