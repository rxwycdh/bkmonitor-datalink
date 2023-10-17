package notifier

import (
	"context"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/window"
	monitorLogger "github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"go.uber.org/zap"
	"sync"
)

type Notifier interface {
	Start()
	Spans() <-chan []window.Span
}

// Options is configuration items for all notifier
type Options struct {
	// Configure for difference queue
	KafkaConfig

	ctx context.Context
	// chanBufferSize The maximum amount of cached data in the queue
	chanBufferSize int
}

type Option func(*Options)

// BufferSize queue chan size
func BufferSize(s int) Option {
	return func(args *Options) {
		args.chanBufferSize = s
	}
}

func Context(ctx context.Context) Option {
	return func(options *Options) {
		options.ctx = ctx
	}
}

type NotifyForm int

const (
	KafkaNotifier NotifyForm = 1 << iota
)

func NewNotifier(form NotifyForm, options ...Option) Notifier {

	switch form {
	case KafkaNotifier:
		return newKafkaNotifier(options...)
	default:
		return newEmptyNotifier()
	}

}

// An emptyNotifier for use when not specified
var (
	once                  sync.Once
	emptyNotifierInstance Notifier
)

type emptyNotifier struct{}

func (e emptyNotifier) Spans() <-chan []window.Span {
	return make(chan []window.Span, 0)
}

func (e emptyNotifier) Start() {}

func newEmptyNotifier() Notifier {
	once.Do(func() {
		emptyNotifierInstance = emptyNotifier{}
	})

	return emptyNotifierInstance
}

var logger = monitorLogger.With(
	zap.String("location", "notifier"),
)
