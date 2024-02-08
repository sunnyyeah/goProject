package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	queueLen = 50
)

type transaction struct {
	rootSpanId uint64
}

type reporter struct {
	tickIntervalsMs int64
	flushIntervalMs int64
	// Transaction队列
	queue chan *transaction
	close chan struct{}

	wg *sync.WaitGroup

	lastFlushTimeMs int64 // 上次 flush 的时间

	closeFlag int32
}

func newReporter() *reporter {
	rpt := &reporter{}

	rpt.tickIntervalsMs = 1000
	rpt.flushIntervalMs = 1000
	rpt.queue = make(chan *transaction, queueLen)
	rpt.close = make(chan struct{})

	rpt.wg = new(sync.WaitGroup)
	return rpt
}

func (reporter *reporter) run() {
	// 定时器
	ticker := time.NewTicker(time.Millisecond * time.Duration(reporter.tickIntervalsMs))
	defer reporter.wg.Done()
	defer ticker.Stop()
	for {
		select {
		case t := <-reporter.queue: // 有数据进来就会执行一次这里代码
			reporter.append(t)
		case <-reporter.close:
			chanLen := len(reporter.queue)
			fmt.Println("close， chanLen = %d", chanLen)
			for i := 0; i < chanLen; i++ {
				select {
				case t := <-reporter.queue:
					reporter.append(t)
				default:
				}
			}
			reporter.flush()
			return
		case <-ticker.C: // 每过 1s 会执行一次这里的代码
			if time.Now().UnixNano()/1e6-reporter.getLastFlushTime() > reporter.flushIntervalMs {
				reporter.flush()
			}
		}
	}
}

func (reporter *reporter) Report(t *transaction) {
	select {
	case reporter.queue <- t:
		fmt.Printf("send %d>>>>>>>>>>>>>>>>\n", t.rootSpanId)
	default:
		fmt.Println("queue full...")
	}
}

func (reporter *reporter) append(t *transaction) {
	fmt.Printf(">>>>>>>>>>>>>>>> append %d\n", t.rootSpanId)
	//time.Sleep(time.Second * 2)
}

func (reporter *reporter) flush() {
	fmt.Print("flush\n")
	reporter.setLastFlushTime(time.Now().UnixNano() / 1e6)
}

func (reporter *reporter) setLastFlushTime(t int64) {
	atomic.StoreInt64(&reporter.lastFlushTimeMs, t)
}

func (reporter *reporter) getLastFlushTime() int64 {
	return atomic.LoadInt64(&reporter.lastFlushTimeMs)
}

func main() {
	rpt := newReporter()
	rpt.wg.Add(1)
	go rpt.run()

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 10)
		rpt.Report(&transaction{
			rootSpanId: uint64(i),
		})
	}

	rpt.close <- struct{}{}
}
