package subpub

import (
	"context"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscription struct {
	subject  string
	handler  MessageHandler
	stopChan chan struct{}
	pubSub   *PubSub
}

func (s *Subscription) Unsubscribe() {
	s.pubSub.unsubscribe(s.subject, s)
}

type PubSub struct {
	subjects *sync.Map
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewSubPub() SubPub {
	ctx, cancel := context.WithCancel(context.Background())
	return &PubSub{
		subjects: &sync.Map{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (ps *PubSub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	sub := &Subscription{
		subject:  subject,
		handler:  cb,
		stopChan: make(chan struct{}),
		pubSub:   ps,
	}

	subs, _ := ps.subjects.LoadOrStore(subject, &sync.Map{})
	subs.(*sync.Map).Store(sub, struct{}{})

	go func() {
		for {
			select {
			case <-sub.stopChan:
				return
			case <-ps.ctx.Done():
				return
			}
		}
	}()

	return sub, nil
}

func (ps *PubSub) Publish(subject string, msg interface{}) error {
	subs, ok := ps.subjects.Load(subject)
	if !ok {
		return nil
	}

	subs.(*sync.Map).Range(func(key, value interface{}) bool {
		sub := key.(*Subscription)
		go sub.handler(msg)
		return true
	})

	return nil
}

func (ps *PubSub) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		ps.cancel()
		return nil
	}
}

func (ps *PubSub) unsubscribe(subject string, sub *Subscription) {
	subs, ok := ps.subjects.Load(subject)
	if !ok {
		return
	}
	subs.(*sync.Map).Delete(sub)
	close(sub.stopChan)
}
