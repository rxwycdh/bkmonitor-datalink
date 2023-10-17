package window

import (
	"context"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/storage"
	monitorLogger "github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"github.com/cespare/xxhash/v2"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

// DistributiveWindowOptions all configs
type DistributiveWindowOptions struct {
	subWindowSize               int
	watchExpiredInterval        time.Duration
	concurrentProcessCount      int
	concurrentExpirationMaximum int
	//// messageConcurrentListener Set the concurrency number to listen for messages in the queue.
	//messageConcurrentListener int
}

type DistributiveWindowOption func(*DistributiveWindowOptions)

// DistributiveWindowSubSize The number of sub-windows, each of which maintains its own data.
func DistributiveWindowSubSize(maxSize int) DistributiveWindowOption {
	return func(options *DistributiveWindowOptions) {
		options.subWindowSize = maxSize
	}
}

// DistributiveWindowWatchExpiredInterval unit: ms. The duration of check expiration trace in window.
// If value is too small, the concurrent performance may be affected
func DistributiveWindowWatchExpiredInterval(interval int) DistributiveWindowOption {
	return func(options *DistributiveWindowOptions) {
		options.watchExpiredInterval = time.Duration(interval) * time.Millisecond
	}
}

// ConcurrentProcessCount The maximum concurrency.
// For example, concurrentProcessCount is set to 10 and subWindowSize is set to 5,
// then each sub-window can have a maximum of 10 traces running at the same time,
// and a total of 5 * 10 can be processed at the same time.
func ConcurrentProcessCount(c int) DistributiveWindowOption {
	return func(options *DistributiveWindowOptions) {
		options.concurrentProcessCount = c
	}
}

// ConcurrentExpirationMaximum Maximum number of concurrent expirations
func ConcurrentExpirationMaximum(c int) DistributiveWindowOption {
	return func(options *DistributiveWindowOptions) {
		options.concurrentExpirationMaximum = c
	}
}

//func MessageConcurrentListener(c int) DistributiveWindowOption {
//	return func(options *DistributiveWindowOptions) {
//		options.messageConcurrentListener = c
//	}
//}

// DistributiveWindow Parent-child window implementation classes.
// where each independent child window maintains its own data, when a span is to be added to the parent window,
// it is added to a child window through hash. Therefore, the child window is the downstream node of the parent window.
// At the same time, the parent window also acts as an observer
// and periodically traverses the child window to check whether there is an expired trace.
// Parent-Window(observer) -----> Child-Window(trigger)
//
//	                                 ^
//	                Via Hash         |
//	Span   ---------------------------
type DistributiveWindow struct {
	dataId string
	config DistributiveWindowOptions

	subWindows      map[int]*distributiveSubWindow
	observers       map[observer]struct{}
	saveRequestChan chan<- storage.SaveRequest

	ctx    context.Context
	logger monitorLogger.Logger
}

func NewDistributiveWindow(dataId string, ctx context.Context, processor Processor, saveReqChan chan<- storage.SaveRequest, specificOptions ...DistributiveWindowOption) Operator {

	specificConfig := &DistributiveWindowOptions{}
	for _, setter := range specificOptions {
		setter(specificConfig)
	}
	window := &DistributiveWindow{
		dataId:          dataId,
		config:          *specificConfig,
		observers:       make(map[observer]struct{}, specificConfig.subWindowSize),
		ctx:             ctx,
		saveRequestChan: saveReqChan,
	}

	// Register sub-windows Event
	subWindowMapping := make(map[int]*distributiveSubWindow, specificConfig.subWindowSize)
	for i := 0; i < specificConfig.subWindowSize; i++ {
		w := newDistributiveSubWindow(dataId, ctx, i, processor, window.saveRequestChan, specificConfig.concurrentExpirationMaximum)
		subWindowMapping[i] = w
		window.register(w)
	}

	window.subWindows = subWindowMapping
	window.logger = monitorLogger.With(
		zap.String("location", "window"),
		zap.String("dataId", dataId),
	)
	return window
}

func (w *DistributiveWindow) locate(uni string) *distributiveSubWindow {
	hashValue := xxhash.Sum64([]byte(uni))
	a := int(hashValue) % 10
	if a < 0 {
		a = -a
	}
	return w.subWindows[a]
}

func (w *DistributiveWindow) Start(spanChan <-chan []Span, runtimeOpts ...RuntimeConfigOption) {

	for ob, _ := range w.observers {
		ob.assembleRuntimeConfig(runtimeOpts...)
		for i := 0; i < w.config.concurrentProcessCount; i++ {
			go ob.handleNotify()
		}
	}

	go w.startWatch()
	w.logger.Infof("DataId: %s created with %d sub-window, %d concurrentProcessCount", w.dataId, len(w.observers), w.config.concurrentProcessCount)

	go w.Handle(spanChan)
	//for i := 0; i < w.config.messageConcurrentListener; i++ {
	//	go w.Handle(spanChan)
	//}
}

func (w *DistributiveWindow) ReportMetric() map[OperatorMetricKey][]OperatorMetric {

	r := make(map[OperatorMetricKey][]OperatorMetric, 2)

	for _, subWindow := range w.subWindows {

		subWindowDimension := map[string]string{"dataId": w.dataId, "subWindowId": strconv.Itoa(subWindow.id)}
		traceCount := 0
		spanCount := 0

		subWindow.m.Range(func(key, value any) bool {
			traceCount++
			v := value.(*CollectTrace)
			spanCount += len(v.Spans)
			return true
		})

		traceV, traceE := r[TraceCount]
		if traceE {
			traceV = append(traceV, OperatorMetric{Value: traceCount, Dimension: subWindowDimension})
			r[TraceCount] = traceV
		} else {
			r[TraceCount] = []OperatorMetric{{Value: traceCount, Dimension: subWindowDimension}}
		}

		spanV, spanE := r[SpanCount]
		if spanE {
			spanV = append(spanV, OperatorMetric{Value: spanCount, Dimension: subWindowDimension})
			r[SpanCount] = spanV
		} else {
			r[SpanCount] = []OperatorMetric{{Value: spanCount, Dimension: subWindowDimension}}
		}
	}

	return r
}

func (w *DistributiveWindow) Handle(spanChan <-chan []Span) {
	w.logger.Infof("DistributiveWindow handle started.")
loop:
	for {
		select {
		case m := <-spanChan:
			for _, span := range m {
				subWindow := w.locate(span.TraceId)
				subWindow.add(span)
			}
		case <-w.ctx.Done():
			w.logger.Infof("Handle span stopped.")
			break loop
		}
	}
}

func (w *DistributiveWindow) register(o observer) {
	w.observers[o] = struct{}{}
}

func (w *DistributiveWindow) deRegister(o observer) {
	delete(w.observers, o)
}

func (w *DistributiveWindow) startWatch() {
	// todo addMetrics: 监听器方法耗时
	tick := time.NewTicker(w.config.watchExpiredInterval)
	w.logger.Infof("DistributiveWindow watching started. interval: %dms", w.config.watchExpiredInterval.Milliseconds())
	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info("trigger watch stopped.")
			return
		case <-tick.C:
			for ob, _ := range w.observers {
				ob.detectNotify()
			}
		}
	}
}

type distributiveSubWindow struct {
	id              int
	dataId          string
	runtimeStrategy ConfigBaseRuntimeStrategies

	eventChan            chan Event
	processor            Processor
	writeSaveRequestChan chan<- storage.SaveRequest

	m      *sync.Map
	ctx    context.Context
	logger monitorLogger.Logger
}

func newDistributiveSubWindow(dataId string, ctx context.Context, index int, processor Processor, saveReqChan chan<- storage.SaveRequest, concurrentMaximum int) *distributiveSubWindow {
	return &distributiveSubWindow{
		id:                   index,
		dataId:               dataId,
		eventChan:            make(chan Event, concurrentMaximum),
		writeSaveRequestChan: saveReqChan,
		processor:            processor,
		m:                    &sync.Map{},
		ctx:                  ctx,
		logger: monitorLogger.With(
			zap.String("location", "window"),
			zap.String("dataId", dataId),
			zap.String("sub-window-id", strconv.Itoa(index)),
		),
	}
}

func (d *distributiveSubWindow) assembleRuntimeConfig(runtimeOpt ...RuntimeConfigOption) {

	config := RuntimeConfig{}
	for _, setter := range runtimeOpt {
		setter(&config)
	}
	d.runtimeStrategy = *NewRuntimeStrategies(
		config,
		[]ReentrantRuntimeStrategy{ReentrantLogRecord, ReentrantLimitMaxCount, RefreshUpdateTime},
		[]ReentrantRuntimeStrategy{PredicateLimitMaxDuration, PredicateNoDataDuration},
	)
}

func (d *distributiveSubWindow) detectNotify() {

	// todo 检查自己是否有过期的traceId
	expiredKeys := make([]string, 0)

	d.m.Range(func(key, value any) bool {
		v := value.(*CollectTrace)
		if d.runtimeStrategy.predicate(&v.Runtime, *v) {
			expiredKeys = append(expiredKeys, key.(string))
		}
		return true
	})

	if len(expiredKeys) > 0 {
		for _, k := range expiredKeys {
			v, exists := d.m.Load(k)
			if !exists {
				d.logger.Errorf("An expired key[%s] was detected but does not exist in the mapping", string(k))
				continue
			}
			d.m.Delete(k)
			d.eventChan <- Event{v.(*CollectTrace)}
		}
	}
}

func (d *distributiveSubWindow) handleNotify() {
loop:
	for {
		select {
		case e := <-d.eventChan:
			d.processor.PreProcess(d.writeSaveRequestChan, e)
		case <-d.ctx.Done():
			break loop
		}
	}
}

func (d *distributiveSubWindow) add(span Span) {

	value, exist := d.m.Load(span.TraceId)
	standardSpan := span.ToStandardSpan()

	if !exist {
		graph := NewDiGraph()
		graph.AddNode(&Node{StandardSpan: standardSpan})
		rt := d.runtimeStrategy.handleNew()
		d.m.Store(span.TraceId, &CollectTrace{
			TraceId: span.TraceId,
			Spans:   []*StandardSpan{standardSpan},
			Graph:   graph,

			Runtime: rt,
		})
	} else {
		collect := value.(*CollectTrace)
		graph := collect.Graph
		graph.AddNode(&Node{StandardSpan: standardSpan})
		collect.Spans = append(collect.Spans, standardSpan)
		collect.Graph = graph

		d.runtimeStrategy.handleExist(&collect.Runtime, *collect)
		d.m.Store(span.TraceId, collect)
	}
}
