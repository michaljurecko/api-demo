// Package shutdown provides simple LIFO stack, to stop internal services in reverse order of their start.
package shutdown

import (
	"context"
)

type Stack struct {
	cancel context.CancelFunc
	stack  []func()
}

func NewStack(ctx context.Context) *Stack {
	ctx, cancel := context.WithCancel(ctx)
	s := &Stack{cancel: cancel}
	go func() {
		<-ctx.Done()
		s.invoke()
	}()
	return s
}

// Shutdown triggers shutdown callbacks manually.
// Shutdown is also triggered when the context passed to NewStack is canceled.
func (s *Stack) Shutdown() {
	s.cancel()
}

func (s *Stack) OnShutdown(f func()) {
	s.stack = append(s.stack, f)
}

func (s *Stack) invoke() {
	for i := len(s.stack) - 1; i >= 0; i-- {
		s.stack[i]()
	}
}
