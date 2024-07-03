package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	defaultShutdown = NewShutdown()
)

func HookShutdownCallback(name string, f func()) {
	defaultShutdown.HookShutdownCallback(name, f)
}

func WaitForShutdown(sigs ...os.Signal) {
	defaultShutdown.WaitForShutdown(sigs...)
}

func SysDown() <-chan struct{} {
	return defaultShutdown.SysDown()
}

func Context() context.Context {
	return defaultShutdown.Context()
}

// define a struct to manage shutdown
type Shutdown struct {
	rootCtx   context.Context
	cancel    func()
	mutex     sync.Mutex
	callbacks []callback
	sigCh     chan os.Signal
}

type callback struct {
	name string
	f    func()
}

func NewShutdown() *Shutdown {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	return &Shutdown{
		rootCtx:   ctx,
		cancel:    cancel,
		callbacks: make([]callback, 0),
		sigCh:     sigCh,
	}
}

func (s *Shutdown) HookShutdownCallback(name string, f func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.callbacks = append(s.callbacks, callback{
		name: name,
		f:    f,
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
	logrus.Info("shutdown signal received. wait for 1 second to begin shutdown...")
	time.Sleep(time.Second)
	s.shutdown()
	logrus.Info("shutdown completed.")
}

func (s *Shutdown) shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	wg := sync.WaitGroup{}
	for _, f := range s.callbacks {
		wg.Add(1)
		go func(f callback) {
			defer func() {
				logrus.Infof("shutdown callback %s done", f.name)
				wg.Done()
			}()
			logrus.Infof("begin shutdown callback %s", f.name)
			f.f()
		}(f)
	}
	wg.Wait()
}
