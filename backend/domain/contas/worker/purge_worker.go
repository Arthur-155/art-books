//gofi:context contas
package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/joaoprofile/gofi-sdk-go/base/errs"
	"github.com/joaoprofile/gofi-sdk-go/gofi"
	"github.com/joaoprofile/gofi-sdk-go/obs/logging"
)

const (
	DefaultPurgeInterval = time.Hour
	purgeRunTimeout      = 5 * time.Minute
)

// Purger apaga o que venceu em sessões, tokens, pedidos e bloqueios.
type Purger interface {
	PurgeExpired(ctx context.Context) errs.AppError
}

// PurgeWorker roda o expurgo a cada intervalo, como componente do gofi.
type PurgeWorker struct {
	purger   Purger
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	done     chan struct{}
}

func NewPurgeWorker(purger Purger, interval time.Duration) *PurgeWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &PurgeWorker{purger: purger, interval: interval, ctx: ctx, cancel: cancel, done: make(chan struct{})}
}

func (w *PurgeWorker) Name() string      { return "contas purge worker" }
func (w *PurgeWorker) Stage() gofi.Stage { return gofi.StageServer }

func (w *PurgeWorker) Start(_ context.Context, _ *gofi.Runtime) error {
	return nil
}

// Run bloqueia até Stop.
func (w *PurgeWorker) Run() error {
	defer close(w.done)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return nil
		case <-ticker.C:
			w.runOnce()
		}
	}
}

func (w *PurgeWorker) Stop(ctx context.Context) error {
	w.cancel()
	select {
	case <-w.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *PurgeWorker) runOnce() {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(w.ctx), purgeRunTimeout)
	defer cancel()
	if appErr := w.purger.PurgeExpired(ctx); appErr.Exists() {
		logging.FromContext(ctx).Error("contas purge worker: failed", slog.Any("error", appErr))
	}
}
