package dump

import (
	"sync"
	"time"
)

const (
	dumpNetEventUnknow = iota
	dumpNetEventTimetick
	dumpNetEventRecvIncr
	dumpNetEventRecvDecr
	dumpNetEventSendIncr
	dumpNetEventSendDecr
)

var onceRT sync.Once
var dumpCtx *dumpContext

type dumpContext struct {
	dumpFlg              bool
	packetRecv           int64
	packetRecvHandle     int64
	packetSend           int64
	packetSendHandle     int64
	packetRecvHandleRate int64
	packetSendHandleRate int64
	dumpEventChan        chan uint64
	dumpFunc             func(int64, int64, int64, int64)
	reportAddr           string
}

// InitDump interval(time.Second int)
func InitDump(dump bool, interval int, reportAddr string, dumpFunc func(int64, int64, int64, int64)) {
	onceRT.Do(func() {
		dumpCtx = &dumpContext{
			dumpFlg:       dump,
			dumpEventChan: make(chan uint64, 1000),
			dumpFunc:      dumpFunc,
			reportAddr:    reportAddr,
		}
		dumpCtx.loopRuntimeInfo(int64(interval))
	})
}

// NetEventRecvIncr 收到一个网络事件包
func NetEventRecvIncr(eventid int) {
	dumpCtx.dumpEventChan <- dumpNetEventRecvIncr | (uint64(eventid) << 32)
}

// NetEventRecvDecr 处理完成一个网络事件包
func NetEventRecvDecr(eventid int) {
	dumpCtx.dumpEventChan <- dumpNetEventRecvDecr | (uint64(eventid) << 32)
}

// NetEventSendIncr 发送一个网络事件包
func NetEventSendIncr(eventid int) {
	dumpCtx.dumpEventChan <- dumpNetEventSendIncr | (uint64(eventid) << 32)
}

// NetEventSendDecr 得到一个发送的网络事件包回应
func NetEventSendDecr(eventid int) {
	dumpCtx.dumpEventChan <- dumpNetEventSendDecr | (uint64(eventid) << 32)
}

func (d *dumpContext) loopRuntimeInfo(interval int64) {
	go func() {
		for {
			e, _ := <-d.dumpEventChan
			op := 0xFFFFFFFF & e
			// eventid := int(0xFFFFFFFF00000000 & e)
			switch op {
			case dumpNetEventTimetick:
				if d.dumpFlg {
					d.packetRecvHandleRate = d.packetRecvHandle / interval
					d.packetSendHandleRate = d.packetSendHandle / interval
					d.dumpFunc(d.packetRecv, d.packetSend, d.packetRecvHandleRate, d.packetSendHandleRate)
				}
				d.packetRecvHandle = 0
				d.packetSendHandle = 0
			case dumpNetEventRecvIncr:
				d.packetRecv++
			case dumpNetEventRecvDecr:
				d.packetRecv--
				d.packetRecvHandle++
			case dumpNetEventSendIncr:
				d.packetSend++
			case dumpNetEventSendDecr:
				d.packetSend--
				d.packetSendHandle++
			}
		}
	}()

	go func() {
		for {
			<-time.After(time.Second * time.Duration(interval))
			d.dumpEventChan <- dumpNetEventTimetick
		}
	}()

	go func() {
		for {
			<-time.After(time.Second * time.Duration(interval))
			d.report()
		}
	}()
}

func (d *dumpContext) report() {
}
