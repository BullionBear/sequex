package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
)

// define a struct to manage shutdown
type Shutdown struct {
	logger    log.Logger
	rootCtx   context.Context
	cancel    func()
	mutex     sync.Mutex
	callbacks []callback
	sigCh     chan os.Signal
}

type callback struct {
	name    string
	f       func()
	timeout time.Duration // not used yet
}

func NewShutdown(log log.Logger) *Shutdown {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	return &Shutdown{
		logger:    log,
		rootCtx:   ctx,
		cancel:    cancel,
		callbacks: make([]callback, 0),
		sigCh:     sigCh,
	}
}

// HookShutdownCallback registers a callback function to be executed during shutdown.
// The timeout parameter specifies how long to wait for the callback to complete.
// If timeout is 0, the callback will run without a timeout.
// If timeout is > 0 and the callback doesn't complete within that time, it will be logged as a timeout error.
func (s *Shutdown) HookShutdownCallback(name string, f func(), timeout time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.callbacks = append(s.callbacks, callback{
		name:    name,
		f:       f,
		timeout: timeout,
	})
}

func (s *Shutdown) Context() context.Context {
	return s.rootCtx
}

func (s *Shutdown) SysDown() <-chan struct{} {
	return s.rootCtx.Done()
}

func (s *Shutdown) WaitForShutdown(sigs ...os.Signal) {
	if len(sigs) > 0 {
		signal.Notify(s.sigCh, sigs...)
	}
	<-s.sigCh
	s.cancel()
	s.logger.Info("shutdown signal received. wait for 1 second to begin shutdown...")
	time.Sleep(time.Second)
	s.shutdown()
	s.logger.Info("shutdown completed.")
}

// ShutdownNow manually triggers the shutdown process.
// This is useful for programmatic shutdown without waiting for system signals.
func (s *Shutdown) ShutdownNow() {
	s.cancel()
	s.logger.Info("manual shutdown triggered. wait for 1 second to begin shutdown...")
	time.Sleep(time.Second)
	s.shutdown()
	s.logger.Info("shutdown completed.")
}

func (s *Shutdown) shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	wg := sync.WaitGroup{}
	for _, f := range s.callbacks {
		wg.Add(1)
		go func(f callback) {
			defer func() {
				wg.Done()
			}()
			s.logger.Info("begin shutdown callback %s", log.String("name", f.name))

			// Create context with timeout if specified
			var ctx context.Context
			var cancel context.CancelFunc
			if f.timeout > 0 {
				ctx, cancel = context.WithTimeout(context.Background(), f.timeout)
				defer cancel()
			} else {
				ctx = context.Background()
			}

			// Execute callback with timeout handling
			done := make(chan struct{})
			go func() {
				defer close(done)
				f.f()
			}()

			select {
			case <-done:
				s.logger.Info("shutdown callback %s done", log.String("name", f.name))
			case <-ctx.Done():
				if f.timeout > 0 {
					s.logger.Error("shutdown callback %s timeout after %v",
						log.String("name", f.name),
						log.Any("timeout", f.timeout))
				}
			}
		}(f)
	}
	wg.Wait()
}
